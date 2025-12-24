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
type generateRequest struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

type generateResponse struct {
	Result []string `json:"result"`
}

type errorResponse struct {
	Error   string   `json:"error"`
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

// Generate handles POST /fizzbuzz - generates a fizzbuzz sequence
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
