package http

import (
	"fizzbuzz-service/internal/infrastructure/http/handler"
	"fizzbuzz-service/internal/infrastructure/http/middleware"
	"log/slog"
	"net/http"
)

func NewRouter(
	fizzBuzzHandler *handler.FizzBuzzHandler,
	statsHandler *handler.StatisticsHandler,
	healthHandler *handler.HealthHandler,
	logger *slog.Logger,

) http.Handler {
	mux := http.NewServeMux()

	// Each handler registers its own routes
	fizzBuzzHandler.RegisterRoutes(mux)
	statsHandler.RegisterRoutes(mux)
	healthHandler.RegisterRoutes(mux)

	// Apply middlewares (order matters: recovery should be outermost)
	var h http.Handler = mux
	h = middleware.LoggingMiddleware(logger)(h)
	h = middleware.RecoveryMiddleware(logger)(h)

	return h
}
