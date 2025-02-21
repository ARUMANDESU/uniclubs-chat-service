package httpapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
	"net/http"
)

type Server struct {
	log *slog.Logger

	HTTPServer *http.Server
}

// New initializes and returns a new instance of the Server struct.
// This function is responsible for setting up the HTTP server with the provided
// configuration and handler. It configures various parameters of the HTTP server
// such as the address to listen on, request handling logic, timeouts, and more.
//
// Parameters:
//   - cfg: A pointer to the config.Config struct containing configuration
//     settings like server address and timeout values.
//   - handler: An http.Handler which handles HTTP requests received by the server.
//     This is typically a router or a middleware chain.
//
// Returns:
//   - A pointer to an initialized Server struct containing the configured http.Server.
//
// Usage:
//
//	This function is usually called during the application's initialization phase
//	to set up the main HTTP server based on the specified configurations.
func New(cfg config.Config, log *slog.Logger, handler http.Handler) Server {
	httpServer := &http.Server{
		Addr:           cfg.HTTP.Address,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    cfg.HTTP.Timeout,
		WriteTimeout:   cfg.HTTP.Timeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
	}

	return Server{
		log:        log,
		HTTPServer: httpServer,
	}
}

// Start starts http server.
//
//   - Panics if the server fails to start.
func (s Server) Start(_ context.Context, callback func(error)) {
	const op = "app.http.start"
	log := s.log.With(slog.String("op", op))

	go func() {
		log.Info("http server is running", slog.String("addr", s.HTTPServer.Addr))
		err := s.HTTPServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start http server", logger.Err(err))
			callback(err)
		}
	}()

}

// Stop gracefully shuts down the server without interrupting any active connections.
// It waits for all the active requests to complete and then shuts down the server.
// This method is typically used for gracefully shutting down the server,
// for instance, when the application is receiving a termination signal.
//
// Parameters:
//   - ctx: A context.Context used to provide a deadline for the shutdown process.
//     The server will wait for active requests to finish until the context deadline.
//
// Returns:
//   - An error if the shutdown process encounters any issues; otherwise, nil.
func (s Server) Stop(ctx context.Context) error {
	const op = "app.http.stop"
	log := s.log.With(slog.String("op", op))

	log.Info("stopping http server")
	err := s.HTTPServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
