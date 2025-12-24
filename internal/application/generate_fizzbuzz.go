package application

import (
	"context"
	"fizzbuzz-service/internal/domain"
	"fizzbuzz-service/internal/domain/entity"
	"fizzbuzz-service/internal/domain/service"
	"log/slog"
)

type GenerateFizzBuzzUseCase struct {
	generator *service.FizzBuzzGenerator
	maxLimit  int
	logger    *slog.Logger
}

func NewGenerateFizzBuzzUseCase(
	generator *service.FizzBuzzGenerator,
	maxLimit int,
	logger *slog.Logger,
) *GenerateFizzBuzzUseCase {
	return &GenerateFizzBuzzUseCase{
		generator: generator,
		maxLimit:  maxLimit,
		logger:    logger,
	}
}

// Execute validates input and generates the sequence
func (uc *GenerateFizzBuzzUseCase) Generate(
	ctx context.Context,
	query entity.FizzBuzzQuery,
) ([]string, error) {
	// Validate with detailed error messages
	validation := query.Validate(uc.maxLimit)
	if !validation.Valid {
		return nil, domain.NewValidationError("invalid parameters", validation.Errors...)
	}

	return uc.generator.Generate(query), nil
}
