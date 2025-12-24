package handler

import (
	"encoding/json"
	"fizzbuzz-service/internal/application"
	"fizzbuzz-service/internal/domain"
	"fizzbuzz-service/internal/domain/entity"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FizzBuzzHandler struct {
	generateUseCase *application.GenerateFizzBuzzUseCase
	logger          *slog.Logger
}

// Request/Response DTOs

// swagger:parameters generateFizzBuzz
type generateFizzBuzzParams struct {
	// FizzBuzz generation parameters
	// in: body
	// required: true
	Body generateRequest
}

// GenerateRequest represents the input for FizzBuzz generation
// swagger:model
type generateRequest struct {
	// First divisor (must be > 0)
	// required: true
	// example: 3
	Int1 int `json:"int1"`
	// Second divisor (must be > 0)
	// required: true
	// example: 5
	Int2 int `json:"int2"`
	// Upper limit for the sequence (must be > 0 and <= MAX_LIMIT)
	// required: true
	// example: 15
	Limit int `json:"limit"`
	// String to replace multiples of int1
	// required: true
	// example: fizz
	Str1 string `json:"str1"`
	// String to replace multiples of int2
	// required: true
	// example: buzz
	Str2 string `json:"str2"`
}

// GenerateResponse contains the FizzBuzz sequence
// swagger:model
type generateResponse struct {
	// The generated FizzBuzz sequence
	// required: true
	// example: ["1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz"]
	Result []string `json:"result"`
}

// ErrorResponse represents an error response
// swagger:model
type errorResponse struct {
	// Error message
	// required: true
	// example: invalid parameters
	Error string `json:"error"`
	// Detailed error messages
	// required: false
	Details []string `json:"details,omitempty"`
}

// NewFizzBuzzHandler creates a new FizzBuzz HTTP handler
func NewFizzBuzzHandler(
	generateUseCase *application.GenerateFizzBuzzUseCase,
	logger *slog.Logger,
) *FizzBuzzHandler {
	return &FizzBuzzHandler{
		generateUseCase: generateUseCase,
		logger:          logger,
	}
}

// RegisterRoutes registers all fizzbuzz-related routes
func (h *FizzBuzzHandler) RegisterRoutes(r chi.Router) {
	r.Post("/fizzbuzz", h.Generate)
}

// swagger:route POST /fizzbuzz fizzbuzz generateFizzBuzz
//
// # Generate FizzBuzz Sequence
//
// Generates a customizable FizzBuzz sequence based on the provided parameters.
// The algorithm replaces numbers divisible by int1 with str1, numbers divisible
// by int2 with str2, and numbers divisible by both with str1+str2.
//
// Responses:
//
//	200: generateResponse
//	400: errorResponse
//	500: errorResponse
func (h *FizzBuzzHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Debug("failed to decode request", "error", err)
		h.writeError(w, http.StatusBadRequest, "invalid JSON body", nil)
		return
	}

	query := entity.FizzBuzzQuery{
		FirstDivisor:  req.Int1,
		SecondDivisor: req.Int2,
		UpperLimit:    req.Limit,
		FirstString:   req.Str1,
		SecondString:  req.Str2,
	}

	result, err := h.generateUseCase.Generate(r.Context(), query)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, generateResponse{Result: result})
}

// swagger:response generateResponse
type generateResponseWrapper struct {
	// in: body
	Body generateResponse
}

// swagger:response errorResponse
type errorResponseWrapper struct {
	// in: body
	Body errorResponse
}

// handleError maps domain errors to HTTP responses
func (h *FizzBuzzHandler) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case domain.ValidationError:
		h.writeError(w, http.StatusBadRequest, e.Message, e.Details)
	default:
		h.logger.Error("unexpected error", "error", err)
		h.writeError(w, http.StatusInternalServerError, "internal server error", nil)
	}
}

func (h *FizzBuzzHandler) writeError(w http.ResponseWriter, status int, message string, details []string) {
	h.writeJSON(w, status, errorResponse{
		Error:   message,
		Details: details,
	})
}

func (h *FizzBuzzHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}
