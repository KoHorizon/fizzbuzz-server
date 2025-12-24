package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fizzbuzz-service/internal/application"
	"fizzbuzz-service/internal/domain/service"
	"fizzbuzz-service/internal/infrastructure/http/handler"
	"fizzbuzz-service/internal/infrastructure/persistence/inmemory"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestFizzBuzzHandler_Integration(t *testing.T) {
	// Setup real dependencies (except external services)
	statsRepo := inmemory.NewStatisticsRepository()
	generator := service.NewFizzBuzzGenerator()
	logger := newTestLogger()
	useCase := application.NewGenerateFizzBuzzUseCase(generator, statsRepo, 10000, logger)
	fizzHandler := handler.NewFizzBuzzHandler(useCase, logger)

	// Create a mux and register routes for proper method routing
	mux := http.NewServeMux()
	fizzHandler.RegisterRoutes(mux)

	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedStatus int
		validateBody   func(t *testing.T, body []byte)
	}{
		{
			name:   "valid request returns fizzbuzz sequence",
			method: http.MethodPost,
			body: map[string]interface{}{
				"int1": 3, "int2": 5, "limit": 15,
				"str1": "fizz", "str2": "buzz",
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)

				result, ok := resp["result"].([]interface{})
				if !ok {
					t.Fatal("missing 'result' field")
				}

				if len(result) != 15 {
					t.Errorf("expected 15 items, got %d", len(result))
				}

				// Check fizzbuzz at position 15
				if result[14] != "fizzbuzz" {
					t.Errorf("expected 'fizzbuzz' at position 15, got %v", result[14])
				}
			},
		},
		{
			name:           "invalid JSON returns 400",
			method:         http.MethodPost,
			body:           "not json",
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)

				if _, ok := resp["error"]; !ok {
					t.Error("expected 'error' field in response")
				}
			},
		},
		{
			name:   "validation error returns detailed errors",
			method: http.MethodPost,
			body: map[string]interface{}{
				"int1": 0, "int2": 0, "limit": 0,
				"str1": "", "str2": "",
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)

				details, ok := resp["details"].([]interface{})
				if !ok {
					t.Fatal("expected 'details' array in response")
				}

				if len(details) != 5 {
					t.Errorf("expected 5 validation errors, got %d", len(details))
				}
			},
		},
		{
			name:   "limit too high returns validation error",
			method: http.MethodPost,
			body: map[string]interface{}{
				"int1": 3, "int2": 5, "limit": 20000,
				"str1": "fizz", "str2": "buzz",
			},
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				json.Unmarshal(body, &resp)

				details := resp["details"].([]interface{})
				found := false
				for _, d := range details {
					if str, ok := d.(string); ok && strings.Contains(str, "limit") {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected error about limit")
				}
			},
		},
		{
			name:           "GET method not allowed",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT method not allowed",
			method:         http.MethodPut,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if tt.body != nil {
				switch v := tt.body.(type) {
				case string:
					reqBody = []byte(v)
				default:
					reqBody, _ = json.Marshal(tt.body)
				}
			}

			req := httptest.NewRequest(tt.method, "/fizzbuzz", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				t.Logf("Response body: %s", w.Body.String())
			}

			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestStatisticsHandler_Integration(t *testing.T) {
	statsRepo := inmemory.NewStatisticsRepository()
	logger := newTestLogger()
	useCase := application.NewGetStatisticsUseCase(statsRepo)
	statsHandler := handler.NewStatisticsHandler(useCase, logger)

	// Create a mux and register routes for proper method routing
	mux := http.NewServeMux()
	statsHandler.RegisterRoutes(mux)

	t.Run("returns empty stats initially", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/statistics", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp map[string]interface{}
		json.NewDecoder(w.Body).Decode(&resp)

		if resp["most_frequent_request"] != nil {
			t.Error("expected null most_frequent_request")
		}

		hits := resp["hits"].(float64)
		if hits != 0 {
			t.Errorf("expected 0 hits, got %f", hits)
		}
	})

	t.Run("POST method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/statistics", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected 405, got %d", w.Code)
		}
	})
}
