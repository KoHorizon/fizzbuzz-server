package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server manages HTTP server lifecycle
type Server struct {
	config Config
	server *http.Server
	logger *slog.Logger
}

// New creates a server instance with production-ready timeouts
func New(config Config, handler http.Handler, logger *slog.Logger) *Server {
	return &Server{
		config: config,
		server: &http.Server{
			Addr:    ":" + config.Port,
			Handler: handler,
			// Production timeouts to prevent resource exhaustion
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
		},
		logger: logger,
	}
}

// Run starts the server and blocks until shutdown signal
func (s *Server) Run() error {
	// Channel for server errors
	serverErr := make(chan error, 1)

	// Start server in goroutine
	go func() {
		s.logger.Info("server starting", "port", s.config.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Wait for shutdown signal or server error
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server failed: %w", err)
	case sig := <-stop:
		s.logger.Info("shutdown signal received", "signal", sig.String())
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), s.config.StopTimeout)
	defer cancel()

	s.logger.Info("shutting down server", "timeout", s.config.StopTimeout)

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}

	s.logger.Info("server stopped gracefully")
	return nil
}
