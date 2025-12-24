package application

import (
	"context"
	"fizzbuzz-service/internal/domain"
	"fizzbuzz-service/internal/domain/entity"
	"fizzbuzz-service/internal/domain/service"
	"log/slog"
	"time"
)

type GenerateFizzBuzzUseCase struct {
	generator    *service.FizzBuzzGenerator
	statsUpdater StatisticsUpdater
	maxLimit     int
	logger       *slog.Logger
}

// StatisticsUpdater is a port for updating statistics
type StatisticsUpdater interface {
	UpdateStats(ctx context.Context, query entity.FizzBuzzQuery) error
}

// NewGenerateFizzBuzzUseCase creates the use case
func NewGenerateFizzBuzzUseCase(
	generator *service.FizzBuzzGenerator,
	statsUpdater StatisticsUpdater,
	maxLimit int,
	logger *slog.Logger,
) *GenerateFizzBuzzUseCase {
	return &GenerateFizzBuzzUseCase{
		generator:    generator,
		statsUpdater: statsUpdater,
		maxLimit:     maxLimit,
		logger:       logger,
	}
}

// Generate validates input and generates the sequence
func (uc *GenerateFizzBuzzUseCase) Generate(
	ctx context.Context,
	query entity.FizzBuzzQuery,
) ([]string, error) {
	// Validate with detailed error messages
	validation := query.Validate(uc.maxLimit)
	if !validation.Valid {
		return nil, domain.NewValidationError("invalid parameters", validation.Errors...)
	}

	// Update statistics asynchronously
	// We use a separate goroutine to not block the main request
	// Errors are logged but don't fail the main request (stats are non-critical)
	go func() {
		// Create a new context with timeout, not tied to the request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := uc.statsUpdater.UpdateStats(ctx, query); err != nil {
			uc.logger.Error("failed to update statistics",
				"error", err,
				"query_key", query.Key(),
			)
		}
	}()

	return uc.generator.Generate(query), nil
}
