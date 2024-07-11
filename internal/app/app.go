package app

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app/httpapp"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/handlers"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/ws"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
	"net/http"
)

type App struct {
	log     *slog.Logger
	cfg     config.Config
	HTTPSvr httpapp.Server
}

func New(cfg config.Config, log *slog.Logger) *App {
	commentService := commentservice.NewComment(commentservice.Config{
		Logger:   log,
		Provider: nil,
		Creator:  nil,
		Updater:  nil,
		Deleter:  nil,
	})

	wsManager, err := ws.NewManager(log, commentService)
	if err != nil {
		log.Error("failed to create websocket manager", logger.Err(err))
		panic(err)
	}

	handler := handlers.NewHandler(log, wsManager)
	handler.RegisterRoutes()

	httpServer := httpapp.New(cfg, handler.Mux)

	return &App{
		log:     log,
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
