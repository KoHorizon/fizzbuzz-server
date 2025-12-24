package http

import (
	"fizzbuzz-service/internal/infrastructure/http/handler"
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
	return h
}
