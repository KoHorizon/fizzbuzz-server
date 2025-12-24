package inmemory

import (
	"context"
	"sync"
	"time"

	"fizzbuzz-service/internal/domain/entity"
)

// StatisticsRepository implements both StatisticsUpdater and StatisticsRepository interfaces
type StatisticsRepository struct {
	mu    sync.RWMutex
	stats map[string]*countEntry
}

type countEntry struct {
	query     entity.FizzBuzzQuery
	hitCount  int64
	lastHitAt time.Time
}

// NewStatisticsRepository creates a thread-safe in-memory repository
func NewStatisticsRepository() *StatisticsRepository {
	return &StatisticsRepository{
		stats: make(map[string]*countEntry),
	}
}

// UpdateStats increments the count for a query pattern
// The key includes ALL parameters (including limit) to correctly track unique requests
func (r *StatisticsRepository) UpdateStats(ctx context.Context, query entity.FizzBuzzQuery) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Use the entity's Key() method for consistent key generation
	key := query.Key()

	if entry, exists := r.stats[key]; exists {
		entry.hitCount++
		entry.lastHitAt = time.Now()
	} else {
		r.stats[key] = &countEntry{
			query:     query,
			hitCount:  1,
			lastHitAt: time.Now(),
		}
	}

	return nil
}

// GetMostFrequent returns the query with the highest hit count
func (r *StatisticsRepository) GetMostFrequent(ctx context.Context) (*entity.StatisticsSummary, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return null for most_frequent_request when no requests have been made
	if len(r.stats) == 0 {
		return &entity.StatisticsSummary{
			MostFrequentQuery: nil,
			HitCount:          0,
		}, nil
	}

	var maxEntry *countEntry
	for _, entry := range r.stats {
		if maxEntry == nil || entry.hitCount > maxEntry.hitCount {
			maxEntry = entry
		}
	}

	return &entity.StatisticsSummary{
		MostFrequentQuery: maxEntry.query.ToResponse(),
		HitCount:          maxEntry.hitCount,
	}, nil
}

// GetStats returns all statistics (useful for debugging/testing)
func (r *StatisticsRepository) GetStats() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]int64, len(r.stats))
	for key, entry := range r.stats {
		result[key] = entry.hitCount
	}
	return result
}

// Clear resets all statistics (useful for testing)
func (r *StatisticsRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stats = make(map[string]*countEntry)
}
