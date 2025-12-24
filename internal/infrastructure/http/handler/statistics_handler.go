package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"fizzbuzz-service/internal/application"
	"fizzbuzz-service/internal/domain/entity"

	"github.com/go-chi/chi/v5"
)

// StatisticsHandler handles HTTP requests for statistics operations
type StatisticsHandler struct {
	getStatsUseCase *application.GetStatisticsUseCase
	logger          *slog.Logger
}

// StatisticsSummary for swagger documentation
// swagger:model StatisticsSummary
type statisticsSummaryResponse struct {
	// The most frequently requested FizzBuzz configuration
	// required: false
	MostFrequentRequest *entity.FizzBuzzQueryResponse `json:"most_frequent_request"`
	// Number of times the most frequent request was made
	// required: true
	// example: 42
	Hits int64 `json:"hits"`
}

// NewStatisticsHandler creates a new Statistics HTTP handler
func NewStatisticsHandler(
	getStatsUseCase *application.GetStatisticsUseCase,
	logger *slog.Logger,
) *StatisticsHandler {
	return &StatisticsHandler{
		getStatsUseCase: getStatsUseCase,
		logger:          logger,
	}
}

// RegisterRoutes registers all statistics-related routes
func (h *StatisticsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/statistics", h.GetMostFrequent)
}

// swagger:route GET /statistics statistics getStatistics
//
// # Get Most Frequent Request
//
// Returns the most frequently requested FizzBuzz configuration and its hit count.
// If no requests have been made yet, returns null for most_frequent_request and 0 hits.
//
// Responses:
//
//	200: statisticsResponse
//	500: errorResponse
func (h *StatisticsHandler) GetMostFrequent(w http.ResponseWriter, r *http.Request) {
	stats, err := h.getStatsUseCase.Get(r.Context())
	if err != nil {
		h.logger.Error("failed to get statistics", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	h.writeJSON(w, http.StatusOK, stats)
}

// swagger:response statisticsResponse
type statisticsResponseWrapper struct {
	// in: body
	Body entity.StatisticsSummary
}

func (h *StatisticsHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}
