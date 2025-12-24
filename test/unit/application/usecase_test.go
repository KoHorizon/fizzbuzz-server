package application_test

import (
	"context"
	"errors"
	"fizzbuzz-service/internal/application"
	"fizzbuzz-service/internal/domain"
	"fizzbuzz-service/internal/domain/entity"
	"fizzbuzz-service/internal/domain/service"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockStatsUpdater records calls for verification
type mockStatsUpdater struct {
	mu        sync.Mutex
	calls     []entity.FizzBuzzQuery
	shouldErr bool
}

func (m *mockStatsUpdater) UpdateStats(ctx context.Context, query entity.FizzBuzzQuery) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldErr {
		return errors.New("mock error")
	}

	m.calls = append(m.calls, query)
	return nil
}

func (m *mockStatsUpdater) getCalls() []entity.FizzBuzzQuery {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]entity.FizzBuzzQuery, len(m.calls))
	copy(result, m.calls)
	return result
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestGenerateFizzBuzzUseCase_Execute(t *testing.T) {
	generator := service.NewFizzBuzzGenerator()
	logger := newTestLogger()

	t.Run("valid request returns correct result", func(t *testing.T) {
		mockUpdater := &mockStatsUpdater{}
		useCase := application.NewGenerateFizzBuzzUseCase(generator, mockUpdater, 100, logger)

		query := entity.FizzBuzzQuery{
			FirstDivisor:  3,
			SecondDivisor: 5,
			UpperLimit:    15,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		result, err := useCase.Generate(context.Background(), query)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 15 {
			t.Errorf("expected 15 results, got %d", len(result))
		}

		// Verify last element is "fizzbuzz" (15 is divisible by both 3 and 5)
		if result[14] != "fizzbuzz" {
			t.Errorf("expected last element to be 'fizzbuzz', got %q", result[14])
		}
	})

	t.Run("returns validation error for zero divisor", func(t *testing.T) {
		mockUpdater := &mockStatsUpdater{}
		useCase := application.NewGenerateFizzBuzzUseCase(generator, mockUpdater, 100, logger)

		query := entity.FizzBuzzQuery{
			FirstDivisor:  0,
			SecondDivisor: 5,
			UpperLimit:    15,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		_, err := useCase.Generate(context.Background(), query)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var validationErr domain.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}
	})

	t.Run("returns validation error when limit exceeds max", func(t *testing.T) {
		mockUpdater := &mockStatsUpdater{}
		useCase := application.NewGenerateFizzBuzzUseCase(generator, mockUpdater, 100, logger)

		query := entity.FizzBuzzQuery{
			FirstDivisor:  3,
			SecondDivisor: 5,
			UpperLimit:    200,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		_, err := useCase.Generate(context.Background(), query)

		if err == nil {
			t.Fatal("expected validation error, got nil")
		}

		var validationErr domain.ValidationError
		if !errors.As(err, &validationErr) {
			t.Errorf("expected ValidationError, got %T", err)
		}

		// Check that error details mention the limit
		found := false
		for _, detail := range validationErr.Details {
			if strings.Contains(detail, "limit") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected error details to mention 'limit', got: %v", validationErr.Details)
		}
	})

	t.Run("returns multiple validation errors", func(t *testing.T) {
		mockUpdater := &mockStatsUpdater{}
		useCase := application.NewGenerateFizzBuzzUseCase(generator, mockUpdater, 100, logger)

		query := entity.FizzBuzzQuery{
			FirstDivisor:  0,
			SecondDivisor: 0,
			UpperLimit:    0,
			FirstString:   "",
			SecondString:  "",
		}

		_, err := useCase.Generate(context.Background(), query)

		var validationErr domain.ValidationError
		if !errors.As(err, &validationErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		if len(validationErr.Details) != 5 {
			t.Errorf("expected 5 validation errors, got %d: %v", len(validationErr.Details), validationErr.Details)
		}
	})

	t.Run("stats updater is called asynchronously", func(t *testing.T) {
		mockUpdater := &mockStatsUpdater{}
		useCase := application.NewGenerateFizzBuzzUseCase(generator, mockUpdater, 100, logger)

		query := entity.FizzBuzzQuery{
			FirstDivisor:  3,
			SecondDivisor: 5,
			UpperLimit:    10,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		_, err := useCase.Generate(context.Background(), query)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Wait for async goroutine
		time.Sleep(50 * time.Millisecond)

		calls := mockUpdater.getCalls()
		if len(calls) != 1 {
			t.Errorf("expected 1 stats update call, got %d", len(calls))
		}

		if len(calls) > 0 && calls[0].UpperLimit != 10 {
			t.Errorf("expected query with limit 10, got %d", calls[0].UpperLimit)
		}
	})

	t.Run("stats error does not fail the request", func(t *testing.T) {
		mockUpdater := &mockStatsUpdater{shouldErr: true}
		useCase := application.NewGenerateFizzBuzzUseCase(generator, mockUpdater, 100, logger)

		query := entity.FizzBuzzQuery{
			FirstDivisor:  3,
			SecondDivisor: 5,
			UpperLimit:    10,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		result, err := useCase.Generate(context.Background(), query)

		// Request should succeed even if stats update fails
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 10 {
			t.Errorf("expected 10 results, got %d", len(result))
		}
	})
}

func TestGetStatisticsUseCase_Execute(t *testing.T) {
	t.Run("returns stats from repository", func(t *testing.T) {
		mockRepo := &mockStatsRepository{
			summary: &entity.StatisticsSummary{
				MostFrequentQuery: &entity.FizzBuzzQueryResponse{
					Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz",
				},
				HitCount: 42,
			},
		}

		useCase := application.NewGetStatisticsUseCase(mockRepo)

		stats, err := useCase.Get(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if stats.HitCount != 42 {
			t.Errorf("expected 42 hits, got %d", stats.HitCount)
		}
	})

	t.Run("returns empty stats when no requests", func(t *testing.T) {
		mockRepo := &mockStatsRepository{
			summary: &entity.StatisticsSummary{
				MostFrequentQuery: nil,
				HitCount:          0,
			},
		}

		useCase := application.NewGetStatisticsUseCase(mockRepo)

		stats, err := useCase.Get(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if stats.HitCount != 0 {
			t.Errorf("expected 0 hits, got %d", stats.HitCount)
		}

		if stats.MostFrequentQuery != nil {
			t.Errorf("expected nil query, got %+v", stats.MostFrequentQuery)
		}
	})
}

type mockStatsRepository struct {
	summary *entity.StatisticsSummary
	err     error
}

func (m *mockStatsRepository) GetMostFrequent(ctx context.Context) (*entity.StatisticsSummary, error) {
	return m.summary, m.err
}
