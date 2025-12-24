package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HealthHandler handles HTTP requests for health checks
type HealthHandler struct{}

type healthResponse struct {
	Status string `json:"status"`
}

// NewHealthHandler creates a new Health HTTP handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// RegisterRoutes registers all health-related routes
func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Check)
}

// Check handles GET /health - returns service health status
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{Status: "healthy"})
}
