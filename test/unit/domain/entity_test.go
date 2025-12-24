package domain_test

import (
	"fizzbuzz-service/internal/domain/entity"
	"strings"
	"testing"
)

// Table Driven Test
func TestFizzBuzzQuery_Validate(t *testing.T) {
	tests := []struct {
		name           string
		query          entity.FizzBuzzQuery
		maxLimit       int
		expectValid    bool
		expectErrors   int
		errorSubstring string // Check if this substring appears in any error
	}{
		{
			name: "valid classic fizzbuzz",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    15,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			maxLimit:    10000,
			expectValid: true,
		},
		{
			name: "valid with limit at max",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  2,
				SecondDivisor: 7,
				UpperLimit:    10000,
				FirstString:   "foo",
				SecondString:  "bar",
			},
			maxLimit:    10000,
			expectValid: true,
		},
		{
			name: "invalid - zero first divisor",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  0,
				SecondDivisor: 5,
				UpperLimit:    15,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			maxLimit:       10000,
			expectValid:    false,
			expectErrors:   1,
			errorSubstring: "int1",
		},
		{
			name: "invalid - negative second divisor",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: -5,
				UpperLimit:    15,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			maxLimit:       10000,
			expectValid:    false,
			expectErrors:   1,
			errorSubstring: "int2",
		},
		{
			name: "invalid - limit exceeds max",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    20000,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			maxLimit:       10000,
			expectValid:    false,
			expectErrors:   1,
			errorSubstring: "limit exceeds",
		},
		{
			name: "invalid - zero limit",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    0,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			maxLimit:       10000,
			expectValid:    false,
			expectErrors:   1,
			errorSubstring: "limit must be greater",
		},
		{
			name: "invalid - empty first string",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    15,
				FirstString:   "",
				SecondString:  "buzz",
			},
			maxLimit:       10000,
			expectValid:    false,
			expectErrors:   1,
			errorSubstring: "str1",
		},
		{
			name: "invalid - empty second string",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 5,
				UpperLimit:    15,
				FirstString:   "fizz",
				SecondString:  "",
			},
			maxLimit:       10000,
			expectValid:    false,
			expectErrors:   1,
			errorSubstring: "str2",
		},
		{
			name: "invalid - multiple errors",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  0,
				SecondDivisor: 0,
				UpperLimit:    0,
				FirstString:   "",
				SecondString:  "",
			},
			maxLimit:     10000,
			expectValid:  false,
			expectErrors: 5, // All fields invalid
		},
		{
			name: "valid - same divisors (edge case)",
			query: entity.FizzBuzzQuery{
				FirstDivisor:  3,
				SecondDivisor: 3,
				UpperLimit:    15,
				FirstString:   "fizz",
				SecondString:  "buzz",
			},
			maxLimit:    10000,
			expectValid: true, // This is valid - same divisors are allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.Validate(tt.maxLimit)

			if result.Valid != tt.expectValid {
				t.Errorf("expected Valid=%v, got Valid=%v", tt.expectValid, result.Valid)
			}

			if tt.expectErrors > 0 && len(result.Errors) != tt.expectErrors {
				t.Errorf("expected %d errors, got %d: %v", tt.expectErrors, len(result.Errors), result.Errors)
			}

			if tt.errorSubstring != "" {
				found := false
				for _, err := range result.Errors {
					if strings.Contains(err, tt.errorSubstring) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error containing %q, got: %v", tt.errorSubstring, result.Errors)
				}
			}
		})
	}
}

func TestFizzBuzzQuery_Key(t *testing.T) {
	tests := []struct {
		name    string
		query1  entity.FizzBuzzQuery
		query2  entity.FizzBuzzQuery
		sameKey bool
	}{
		{
			name: "identical queries have same key",
			query1: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
				FirstString: "fizz", SecondString: "buzz",
			},
			query2: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
				FirstString: "fizz", SecondString: "buzz",
			},
			sameKey: true,
		},
		{
			name: "different limits have different keys",
			query1: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
				FirstString: "fizz", SecondString: "buzz",
			},
			query2: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 100,
				FirstString: "fizz", SecondString: "buzz",
			},
			sameKey: false,
		},
		{
			name: "different divisors have different keys",
			query1: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
				FirstString: "fizz", SecondString: "buzz",
			},
			query2: entity.FizzBuzzQuery{
				FirstDivisor: 2, SecondDivisor: 7, UpperLimit: 15,
				FirstString: "fizz", SecondString: "buzz",
			},
			sameKey: false,
		},
		{
			name: "different strings have different keys",
			query1: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
				FirstString: "fizz", SecondString: "buzz",
			},
			query2: entity.FizzBuzzQuery{
				FirstDivisor: 3, SecondDivisor: 5, UpperLimit: 15,
				FirstString: "foo", SecondString: "bar",
			},
			sameKey: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1 := tt.query1.Key()
			key2 := tt.query2.Key()

			if tt.sameKey && key1 != key2 {
				t.Errorf("expected same key, got %q and %q", key1, key2)
			}

			if !tt.sameKey && key1 == key2 {
				t.Errorf("expected different keys, both are %q", key1)
			}
		})
	}
}
