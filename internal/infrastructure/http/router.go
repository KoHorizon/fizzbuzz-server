package http

import (
	"fizzbuzz-service/internal/infrastructure/http/handler"
	"log/slog"
	"net/http"
	"time"

	custommw "fizzbuzz-service/internal/infrastructure/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	fizzBuzzHandler *handler.FizzBuzzHandler,
	statsHandler *handler.StatisticsHandler,
	healthHandler *handler.HealthHandler,
	logger *slog.Logger,

) http.Handler {
	r := chi.NewRouter()

	// Middleware stack (top = outermost, executes first)
	r.Use(middleware.RequestID)                 // Chi: inject X-Request-Id
	r.Use(middleware.RealIP)                    // Chi: get real IP
	r.Use(custommw.RecoveryMiddleware(logger))  // Custom: slog + JSON response
	r.Use(custommw.LoggingMiddleware(logger))   // Custom: slog structured logging
	r.Use(middleware.Timeout(30 * time.Second)) // Chi: request timeout

	// Register routes
	fizzBuzzHandler.RegisterRoutes(r)
	statsHandler.RegisterRoutes(r)
	healthHandler.RegisterRoutes(r)

	return r
}
