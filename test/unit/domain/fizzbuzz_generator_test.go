package domain_test

import (
	"testing"

	"fizzbuzz-service/internal/domain/entity"
	"fizzbuzz-service/internal/domain/service"
)

func TestFizzBuzzGenerator_Generate(t *testing.T) {
	generator := service.NewFizzBuzzGenerator()

	tests := []struct {
		name     string
		query    entity.FizzBuzzQuery
		expected []string
	}{
		{
			name: "classic fizzbuzz up to 15",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    15,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			expected: []string{
				"1", "2", "fizz", "4", "buzz",
				"fizz", "7", "8", "fizz", "buzz",
				"11", "fizz", "13", "14", "fizzbuzz",
			},
		},
		{
			name: "custom divisors and strings",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  2,
				SecondDivisor: 7,
				UpperLimit:    14,
				FirstString:   "two",
				SecondString:  "seven",
			},
			expected: []string{
				"1", "two", "3", "two", "5", "two", "seven",
				"two", "9", "two", "11", "two", "13", "twoseven",
			},
		},
		{
			name: "limit of 1",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    1,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			expected: []string{"1"},
		},
		{
			name: "same divisors - both strings concatenated",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 3,
				UpperLimit:    6,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			expected: []string{"1", "2", "fizzbuzz", "4", "5", "fizzbuzz"},
		},
		{
			name: "divisor of 1 - every number replaced",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  1,
				SecondDivisor: 2,
				UpperLimit:    5,
				FirstString:   "one",
				SecondString:  "two",
			},
			expected: []string{"one", "onetwo", "one", "onetwo", "one"},
		},
		{
			name: "large divisors beyond limit",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  100,
				SecondDivisor: 200,
				UpperLimit:    5,
				FirstString:   "hundred",
				SecondString:  "twohundred",
			},
			expected: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "special characters in strings",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  2,
				SecondDivisor: 3,
				UpperLimit:    6,
				FirstString:   "🎉",
				SecondString:  "✨",
			},
			expected: []string{"1", "🎉", "✨", "🎉", "5", "🎉✨"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.Generate(tt.query)

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d elements, got %d", len(tt.expected), len(result))
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("position %d: expected %q, got %q", i+1, tt.expected[i], v)
				}
			}
		})
	}
}

// Benchmark for performance verification
func BenchmarkFizzBuzzGenerator_Generate(b *testing.B) {
	generator := service.NewFizzBuzzGenerator()
	query := entity.FizzBuzzQuery{
		FirstDivisor:  3,
		SecondDivisor: 5,
		UpperLimit:    10000,
		FirstString:   "fizz",
		SecondString:  "buzz",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator.Generate(query)
	}
}
