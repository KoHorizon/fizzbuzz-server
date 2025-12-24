package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fizzbuzz-service/internal/application"
	"fizzbuzz-service/internal/domain/service"
	infrahttp "fizzbuzz-service/internal/infrastructure/http"
	"fizzbuzz-service/internal/infrastructure/http/handler"
	"fizzbuzz-service/internal/infrastructure/persistence/inmemory"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func setupTestServer(t *testing.T) (string, func()) {
	t.Helper()

	// Setup dependencies
	statsRepo := inmemory.NewStatisticsRepository()
	generator := service.NewFizzBuzzGenerator()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	generateUseCase := application.NewGenerateFizzBuzzUseCase(generator, statsRepo, 10000, logger)
	getStatsUseCase := application.NewGetStatisticsUseCase(statsRepo)

	fizzHandler := handler.NewFizzBuzzHandler(generateUseCase, logger)
	statsHandler := handler.NewStatisticsHandler(getStatsUseCase, logger)
	healthHandler := handler.NewHealthHandler()

	router := infrahttp.NewRouter(fizzHandler, statsHandler, healthHandler, logger)

	// Create listener on random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	server := &http.Server{Handler: router}

	// Start server
	go server.Serve(listener)

	// Return address and cleanup function
	addr := listener.Addr().String()
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}

	// Wait for server to be ready
	time.Sleep(50 * time.Millisecond)

	return addr, cleanup
}

func TestE2E_FizzBuzzEndpoint(t *testing.T) {
	addr, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("generates correct fizzbuzz sequence", func(t *testing.T) {
		body := map[string]interface{}{
			"int1": 3, "int2": 5, "limit": 15,
			"str1": "fizz", "str2": "buzz",
		}

		resp, err := http.Post(
			fmt.Sprintf("http://%s/fizzbuzz", addr),
			"application/json",
			bytes.NewBuffer(mustMarshal(body)),
		)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}

		var result map[string][]string
		json.NewDecoder(resp.Body).Decode(&result)

		expected := []string{
			"1", "2", "fizz", "4", "buzz",
			"fizz", "7", "8", "fizz", "buzz",
			"11", "fizz", "13", "14", "fizzbuzz",
		}

		if len(result["result"]) != len(expected) {
			t.Fatalf("expected %d items, got %d", len(expected), len(result["result"]))
		}

		for i, v := range result["result"] {
			if v != expected[i] {
				t.Errorf("position %d: expected %q, got %q", i+1, expected[i], v)
			}
		}
	})

	t.Run("returns validation errors with details", func(t *testing.T) {
		body := map[string]interface{}{
			"int1": 0, "int2": -1, "limit": 0,
			"str1": "", "str2": "",
		}

		resp, err := http.Post(
			fmt.Sprintf("http://%s/fizzbuzz", addr),
			"application/json",
			bytes.NewBuffer(mustMarshal(body)),
		)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.StatusCode)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		details, ok := result["details"].([]interface{})
		if !ok {
			t.Fatal("expected 'details' array")
		}

		if len(details) < 4 {
			t.Errorf("expected at least 4 validation errors, got %d", len(details))
		}
	})
}

func TestE2E_StatisticsEndpoint(t *testing.T) {
	addr, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("returns empty stats initially", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://%s/statistics", addr))
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		var stats map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&stats)

		if stats["most_frequent_request"] != nil {
			t.Error("expected null most_frequent_request initially")
		}

		hits := stats["hits"].(float64)
		if hits != 0 {
			t.Errorf("expected 0 hits, got %f", hits)
		}
	})
}

func TestE2E_FullFlow(t *testing.T) {
	addr, cleanup := setupTestServer(t)
	defer cleanup()

	// Make multiple concurrent requests
	const numRequests = 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()

			body := map[string]interface{}{
				"int1": 3, "int2": 5, "limit": 10,
				"str1": "fizz", "str2": "buzz",
			}

			resp, err := http.Post(
				fmt.Sprintf("http://%s/fizzbuzz", addr),
				"application/json",
				bytes.NewBuffer(mustMarshal(body)),
			)
			if err != nil {
				t.Errorf("request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected 200, got %d", resp.StatusCode)
			}
		}()
	}

	wg.Wait()

	// Wait for async stats updates
	time.Sleep(200 * time.Millisecond)

	// Verify statistics
	resp, err := http.Get(fmt.Sprintf("http://%s/statistics", addr))
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}
	defer resp.Body.Close()

	var stats map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&stats)

	hits := stats["hits"].(float64)
	if hits < float64(numRequests) {
		t.Errorf("expected at least %d hits, got %f", numRequests, hits)
	}

	// Verify the most frequent request parameters
	mostFrequent := stats["most_frequent_request"].(map[string]interface{})

	if mostFrequent["int1"].(float64) != 3 {
		t.Errorf("expected int1=3, got %v", mostFrequent["int1"])
	}
	if mostFrequent["int2"].(float64) != 5 {
		t.Errorf("expected int2=5, got %v", mostFrequent["int2"])
	}
	if mostFrequent["limit"].(float64) != 10 {
		t.Errorf("expected limit=10, got %v", mostFrequent["limit"])
	}
}

func TestE2E_DifferentQueriesTrackedSeparately(t *testing.T) {
	addr, cleanup := setupTestServer(t)
	defer cleanup()

	// Make requests with different parameters
	queries := []map[string]interface{}{
		{"int1": 3, "int2": 5, "limit": 10, "str1": "fizz", "str2": "buzz"},
		{"int1": 3, "int2": 5, "limit": 10, "str1": "fizz", "str2": "buzz"},
		{"int1": 3, "int2": 5, "limit": 10, "str1": "fizz", "str2": "buzz"},
		{"int1": 2, "int2": 7, "limit": 20, "str1": "foo", "str2": "bar"},
		{"int1": 2, "int2": 7, "limit": 20, "str1": "foo", "str2": "bar"},
	}

	for _, q := range queries {
		resp, err := http.Post(
			fmt.Sprintf("http://%s/fizzbuzz", addr),
			"application/json",
			bytes.NewBuffer(mustMarshal(q)),
		)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		resp.Body.Close()
	}

	// Wait for async stats updates
	time.Sleep(100 * time.Millisecond)

	// Get statistics
	resp, err := http.Get(fmt.Sprintf("http://%s/statistics", addr))
	if err != nil {
		t.Fatalf("failed to get stats: %v", err)
	}
	defer resp.Body.Close()

	var stats map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&stats)

	// The first query pattern (3 times) should be most frequent
	mostFrequent := stats["most_frequent_request"].(map[string]interface{})

	if mostFrequent["int1"].(float64) != 3 {
		t.Errorf("expected most frequent int1=3, got %v", mostFrequent["int1"])
	}

	hits := stats["hits"].(float64)
	if hits != 3 {
		t.Errorf("expected 3 hits for most frequent, got %f", hits)
	}
}

func TestE2E_HealthEndpoint(t *testing.T) {
	addr, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(fmt.Sprintf("http://%s/health", addr))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var health map[string]string
	json.NewDecoder(resp.Body).Decode(&health)

	if health["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %q", health["status"])
	}
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
