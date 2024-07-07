package app

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app/httpapp"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/ws"
	"log/slog"
	"net/http"
)

type App struct {
	log     *slog.Logger
	cfg     config.Config
	HTTPSvr httpapp.Server
}

func New(cfg config.Config, logger *slog.Logger) *App {
	mux := http.NewServeMux()

	manager := ws.NewManager()

	mux.HandleFunc("/ws", manager.ServeWS)

	httpServer := httpapp.New(cfg, mux)

	return &App{
		log:     logger,
		cfg:     cfg,
		HTTPSvr: httpServer,
	}
}

// Stop gracefully stops the app
//
//	 It stops the following services:
//		- gRPC server
//		- AMQP server
//		- MongoDB
//	 It waits for all background works to be completed
func (a *App) Stop() {
	const op = "app.stop"
	log := a.log.With(slog.String("op", op))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if err := a.HTTPSvr.Stop(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("shutdown error", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
}
