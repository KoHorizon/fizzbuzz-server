package service

import (
	"fizzbuzz-service/internal/domain/entity"
	"strconv"
)

// FizzBuzzGenerator implements the core business logic
type FizzBuzzGenerator struct{}

// NewFizzBuzzGenerator creates a new generator
func NewFizzBuzzGenerator() *FizzBuzzGenerator {
	return &FizzBuzzGenerator{}
}

// Generate creates the fizzbuzz sequence
// Precondition: query has been validated
func (g *FizzBuzzGenerator) Generate(query entity.FizzBuzzQuery) []string {
	result := make([]string, 0, query.UpperLimit)

	for i := 1; i <= query.UpperLimit; i++ {
		result = append(result, g.generateSingle(i, query))
	}

	return result
}

// generateSingle determines the output for a single number
func (g *FizzBuzzGenerator) generateSingle(n int, query entity.FizzBuzzQuery) string {
	divisibleByFirst := n%query.FirstDivisor == 0
	divisibleBySecond := n%query.SecondDivisor == 0

	switch {
	case divisibleByFirst && divisibleBySecond:
		return query.FirstString + query.SecondString
	case divisibleByFirst:
		return query.FirstString
	case divisibleBySecond:
		return query.SecondString
	default:
		return strconv.Itoa(n)
	}
}
