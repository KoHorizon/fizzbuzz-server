package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HealthHandler handles HTTP requests for health checks
type HealthHandler struct{}

// HealthResponse represents the health status
// swagger:model
type healthResponse struct {
	// Service health status
	// required: true
	// example: healthy
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

// swagger:route GET /health health healthCheck
//
// # Health Check
//
// Returns the current health status of the service.
// This endpoint is designed for Kubernetes/Docker health checks.
//
// Responses:
//
//	200: healthResponse
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse{Status: "healthy"})
}

// swagger:response healthResponse
type healthResponseWrapper struct {
	// in: body
	Body healthResponse
}
