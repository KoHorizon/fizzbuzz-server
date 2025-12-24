package http

import (
	"fizzbuzz-service/internal/infrastructure/http/handler"
	"fizzbuzz-service/internal/infrastructure/http/middleware"
	"log/slog"
	"net/http"
)

func NewRouter(
	fizzBuzzHandler *handler.FizzBuzzHandler,
	healthHandler *handler.HealthHandler,
	logger *slog.Logger,

) http.Handler {
	mux := http.NewServeMux()
	fizzBuzzHandler.RegisterRoutes(mux)
	healthHandler.RegisterRoutes(mux)

	var h http.Handler = mux
	h = middleware.LoggingMiddleware(logger)(h)
	h = middleware.RecoveryMiddleware(logger)(h)
	return h
}
