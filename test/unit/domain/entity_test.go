package domain_test

import (
	"fizzbuzz-service/domain/entity"
	"testing"
)

func TestFizzBuzzQuery_Validate(t *testing.T) {
	query := entity.FizzBuzzQuery{
		FirstDivisor:  3,
		SecondDivisor: 5,
		UpperLimit:    15,
		FirstString:   "fizz",
		SecondString:  "buzz",
	}

	result := query.Validate(10000)

	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}
