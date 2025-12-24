package inmemory_test

import (
	"context"
	"sync"
	"testing"

	"fizzbuzz-service/internal/domain/entity"
	"fizzbuzz-service/internal/infrastructure/persistence/inmemory"
)

func TestStatisticsRepository_UpdateStats(t *testing.T) {
	t.Run("increments count for same query", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()
		ctx := context.Background()

		query := entity.FizzBuzzQuery{
			FirstDivisor:  3,
			SecondDivisor: 5,
			UpperLimit:    15,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		// Update 5 times
		for i := 0; i < 5; i++ {
			err := repo.UpdateStats(ctx, query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}

		stats, err := repo.GetMostFrequent(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if stats.HitCount != 5 {
			t.Errorf("expected 5 hits, got %d", stats.HitCount)
		}
	})

	t.Run("tracks different queries separately", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()
		ctx := context.Background()

		query1 := entity.FizzBuzzQuery{
			FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
			FirstString: "fizz", SecondString: "buzz",
		}

		query2 := entity.FizzBuzzQuery{
			FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 100, // Different limit
			FirstString: "fizz", SecondString: "buzz",
		}

		// Query1: 3 times
		for i := 0; i < 3; i++ {
			repo.UpdateStats(ctx, query1)
		}

		// Query2: 5 times
		for i := 0; i < 5; i++ {
			repo.UpdateStats(ctx, query2)
		}

		stats, _ := repo.GetMostFrequent(ctx)

		if stats.HitCount != 5 {
			t.Errorf("expected most frequent to have 5 hits, got %d", stats.HitCount)
		}

		if stats.MostFrequentQuery.Limit != 100 {
			t.Errorf("expected most frequent query limit to be 100, got %d", stats.MostFrequentQuery.Limit)
		}
	})

	t.Run("different strings are tracked separately", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()
		ctx := context.Background()

		query1 := entity.FizzBuzzQuery{
			FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
			FirstString: "fizz", SecondString: "buzz",
		}

		query2 := entity.FizzBuzzQuery{
			FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
			FirstString: "foo", SecondString: "bar", // Different strings
		}

		repo.UpdateStats(ctx, query1)
		repo.UpdateStats(ctx, query2)
		repo.UpdateStats(ctx, query2)

		stats, _ := repo.GetMostFrequent(ctx)

		if stats.MostFrequentQuery.Str1 != "foo" {
			t.Errorf("expected most frequent str1 to be 'foo', got %q", stats.MostFrequentQuery.Str1)
		}
	})
}

func TestStatisticsRepository_GetMostFrequent(t *testing.T) {
	t.Run("returns nil query when empty", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()

		stats, err := repo.GetMostFrequent(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if stats.MostFrequentQuery != nil {
			t.Errorf("expected nil query, got %+v", stats.MostFrequentQuery)
		}

		if stats.HitCount != 0 {
			t.Errorf("expected 0 hits, got %d", stats.HitCount)
		}
	})

	t.Run("returns correct JSON format", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()
		ctx := context.Background()

		query := entity.FizzBuzzQuery{
			FirstDivisor:  3,
			SecondDivisor: 5,
			UpperLimit:    15,
			FirstString:   "fizz",
			SecondString:  "buzz",
		}

		repo.UpdateStats(ctx, query)

		stats, _ := repo.GetMostFrequent(ctx)

		// Verify the response format uses API field names
		if stats.MostFrequentQuery.Int1 != 3 {
			t.Errorf("expected Int1=3, got %d", stats.MostFrequentQuery.Int1)
		}
		if stats.MostFrequentQuery.Int2 != 5 {
			t.Errorf("expected Int2=5, got %d", stats.MostFrequentQuery.Int2)
		}
		if stats.MostFrequentQuery.Limit != 15 {
			t.Errorf("expected Limit=15, got %d", stats.MostFrequentQuery.Limit)
		}
	})
}

func TestStatisticsRepository_Concurrency(t *testing.T) {
	t.Run("handles concurrent updates safely", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()
		ctx := context.Background()

		query := entity.FizzBuzzQuery{
			FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
			FirstString: "fizz", SecondString: "buzz",
		}

		const numGoroutines = 100
		const updatesPerGoroutine = 100

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < updatesPerGoroutine; j++ {
					repo.UpdateStats(ctx, query)
				}
			}()
		}

		wg.Wait()

		stats, err := repo.GetMostFrequent(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := int64(numGoroutines * updatesPerGoroutine)
		if stats.HitCount != expected {
			t.Errorf("expected %d hits, got %d", expected, stats.HitCount)
		}
	})

	t.Run("handles concurrent reads and writes", func(t *testing.T) {
		repo := inmemory.NewStatisticsRepository()
		ctx := context.Background()

		query := entity.FizzBuzzQuery{
			FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
			FirstString: "fizz", SecondString: "buzz",
		}

		// Pre-populate
		repo.UpdateStats(ctx, query)

		const numGoroutines = 50
		var wg sync.WaitGroup
		wg.Add(numGoroutines * 2) // Half writers, half readers

		// Writers
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					repo.UpdateStats(ctx, query)
				}
			}()
		}

		// Readers
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					stats, err := repo.GetMostFrequent(ctx)
					if err != nil {
						t.Errorf("read error: %v", err)
					}
					if stats.HitCount < 1 {
						t.Errorf("expected at least 1 hit")
					}
				}
			}()
		}

		wg.Wait()
	})
}

func TestStatisticsRepository_Clear(t *testing.T) {
	repo := inmemory.NewStatisticsRepository()
	ctx := context.Background()

	query := entity.FizzBuzzQuery{
		FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
		FirstString: "fizz", SecondString: "buzz",
	}

	repo.UpdateStats(ctx, query)
	repo.UpdateStats(ctx, query)

	repo.Clear()

	stats, _ := repo.GetMostFrequent(ctx)

	if stats.HitCount != 0 {
		t.Errorf("expected 0 hits after clear, got %d", stats.HitCount)
	}
}
