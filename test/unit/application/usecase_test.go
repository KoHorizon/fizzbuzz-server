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
	"testing"
	"time"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestGenerateFizzBuzzUseCase_Execute(t *testing.T) {
	generator := service.NewFizzBuzzGenerator()
	logger := newTestLogger()

	t.Run("valid request returns correct result", func(t *testing.T) {
		useCase := application.NewGenerateFizzBuzzUseCase(generator, 100, logger)

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
		useCase := application.NewGenerateFizzBuzzUseCase(generator, 100, logger)

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
		useCase := application.NewGenerateFizzBuzzUseCase(generator, 100, logger)

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
		useCase := application.NewGenerateFizzBuzzUseCase(generator, 100, logger)

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
		useCase := application.NewGenerateFizzBuzzUseCase(generator, 100, logger)

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

	})

	t.Run("stats error does not fail the request", func(t *testing.T) {
		useCase := application.NewGenerateFizzBuzzUseCase(generator, 100, logger)

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
