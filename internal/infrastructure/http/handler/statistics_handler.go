package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"fizzbuzz-service/internal/application"
)

// StatisticsHandler handles HTTP requests for statistics operations
type StatisticsHandler struct {
	getStatsUseCase *application.GetStatisticsUseCase
	logger          *slog.Logger
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
func (h *StatisticsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /statistics", h.GetMostFrequent)
}

// GetMostFrequent handles GET /statistics - returns the most frequent request
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

func (h *StatisticsHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}
