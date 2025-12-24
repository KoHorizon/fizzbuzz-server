package entity

import "fmt"

type FizzBuzzQuery struct {
	FirstDivisor  int
	SecondDivisor int
	UpperLimit    int
	FirstString   string
	SecondString  string
}

type ValidationResult struct {
	Valid  bool
	Errors []string
}

func (q *FizzBuzzQuery) Validate(maxLimit int) ValidationResult {
	var errors []string

	if q.FirstDivisor <= 0 {
		errors = append(errors, "int1 must be greater than 0")
	}

	if q.SecondDivisor <= 0 {
		errors = append(errors, "int2 must be greater than 0")
	}

	if q.UpperLimit <= 0 {
		errors = append(errors, "limit must be greater than 0")
	} else if q.UpperLimit > maxLimit {
		errors = append(errors, fmt.Sprintf("limit exceeds maximum allowed value of %d", maxLimit))
	}

	if q.FirstString == "" {
		errors = append(errors, "str1 cannot be empty")
	}

	if q.SecondString == "" {
		errors = append(errors, "str2 cannot be empty")
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

// Key generates a unique identifier for this query (used for statistics)
// Includes ALL parameters to correctly track unique request patterns
func (q FizzBuzzQuery) Key() string {
	return fmt.Sprintf("%d:%d:%d:%s:%s",
		q.FirstDivisor,
		q.SecondDivisor,
		q.UpperLimit,
		q.FirstString,
		q.SecondString,
	)
}
