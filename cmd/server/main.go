package main

import (
	"log/slog"
	"os"

	"fizzbuzz-service/internal/application"
	"fizzbuzz-service/internal/domain/service"
	"fizzbuzz-service/internal/infrastructure/config"
	infrahttp "fizzbuzz-service/internal/infrastructure/http"
	"fizzbuzz-service/internal/infrastructure/http/handler"
	"fizzbuzz-service/internal/infrastructure/server"
)

func main() {
	// 1. Load configuration
	cfg := config.Load()

	// 2. Setup structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(cfg.LogLevel),
	}))

	// 3. Wire dependencies (manual DI - could use wire/fx for larger apps)
	generator := service.NewFizzBuzzGenerator()

	generateUseCase := application.NewGenerateFizzBuzzUseCase(generator, cfg.MaxLimit, logger)

	fizzHandler := handler.NewFizzBuzzHandler(generateUseCase, logger)
	healthHandler := handler.NewHealthHandler()

	router := infrahttp.NewRouter(fizzHandler, healthHandler, logger)

	// 4. Configure and run server
	serverCfg := server.Default()
	serverCfg.Port = cfg.Port

	srv := server.New(serverCfg, router, logger)
	if err := srv.Run(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
