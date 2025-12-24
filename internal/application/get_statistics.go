package application

import (
	"context"

	"fizzbuzz-service/internal/domain/entity"
)

// GetStatisticsUseCase retrieves the most frequent request
type GetStatisticsUseCase struct {
	repo StatisticsRepository
}

// StatisticsRepository is a port for reading statistics
type StatisticsRepository interface {
	GetMostFrequent(ctx context.Context) (*entity.StatisticsSummary, error)
}

// NewGetStatisticsUseCase creates the use case
func NewGetStatisticsUseCase(repo StatisticsRepository) *GetStatisticsUseCase {
	return &GetStatisticsUseCase{repo: repo}
}

// Get returns the most frequent request statistics
func (uc *GetStatisticsUseCase) Get(ctx context.Context) (*entity.StatisticsSummary, error) {
	return uc.repo.GetMostFrequent(ctx)
}
