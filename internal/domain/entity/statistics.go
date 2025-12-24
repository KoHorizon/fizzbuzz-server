package entity

// StatisticsSummary represents the most frequent request
type StatisticsSummary struct {
	MostFrequentQuery *FizzBuzzQueryResponse `json:"most_frequent_request"`
	HitCount          int64                  `json:"hits"`
}

// FizzBuzzQueryResponse is the JSON representation of a query
// Separate from FizzBuzzQuery to control API contract
type FizzBuzzQueryResponse struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// DTO mapper to converts a FizzBuzzQuery to its API response format
func (q FizzBuzzQuery) ToResponse() *FizzBuzzQueryResponse {
	return &FizzBuzzQueryResponse{
		Int1:  q.FirstDivisor,
		Int2:  q.SecondDivisor,
		Limit: q.UpperLimit,
		Str1:  q.FirstString,
		Str2:  q.SecondString,
	}
}
