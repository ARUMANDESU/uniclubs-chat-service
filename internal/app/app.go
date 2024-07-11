package app

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app/httpapp"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/handlers"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage/mongodb"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/ws"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
)

type App struct {
	log      *slog.Logger
	cfg      config.Config
	stoppers []Stopper
	starters []Starter
}

type Starter interface {
	Start(ctx context.Context) error
}

type Stopper interface {
	Stop(ctx context.Context) error
}

// New creates a new App instance
//
//	It initializes all the services and repositories required by the app
//	It creates a new HTTP server instance
//	It returns a new App instance
func New(cfg config.Config, log *slog.Logger) *App {
	const op = "app.new"
	l := log.With(slog.String("op", op))

	starters := make([]Starter, 0)
	stoppers := make([]Stopper, 0)

	mongoStorage, err := mongodb.NewStorage(context.Background(), cfg.MongoDB)
	if err != nil {
		l.Error("failed to create mongodb storage", logger.Err(err))
		panic(err)
	}
	stoppers = append(stoppers, &mongoStorage)

	commentService := commentservice.NewComment(commentservice.Config{
		Logger:   log,
		Provider: &mongoStorage,
		Creator:  &mongoStorage,
		Updater:  &mongoStorage,
		Deleter:  &mongoStorage,
	})

	wsManager, err := ws.NewManager(log, commentService)
	if err != nil {
		l.Error("failed to create websocket manager", logger.Err(err))
		panic(err)
	}
	stoppers = append(stoppers, wsManager)

	handler := handlers.NewHandler(log, wsManager)
	handler.RegisterRoutes()

	httpServer := httpapp.New(cfg, handler.Mux)
	starters = append(starters, httpServer)
	stoppers = append(stoppers, httpServer)

	return &App{
		log:      log,
		cfg:      cfg,
		stoppers: stoppers,
		starters: starters,
	}
}

func (a *App) Start() {
	const op = "app.start"
	log := a.log.With(slog.String("op", op))

	for _, s := range a.starters {
		go func(s Starter) {
			if err := s.Start(context.Background()); err != nil {
				log.Error("failed to start service", logger.Err(err))
			}
		}(s)
	}

}

// Stop gracefully stops the app
//
//	It stops services
//	It waits for all background works to be completed
func (a *App) Stop() {
	const op = "app.stop"
	log := a.log.With(slog.String("op", op))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	for _, s := range a.stoppers {
		go func(s Stopper) {
			if err := s.Stop(shutdownCtx); err != nil {
				log.Error("failed to stop service", logger.Err(err))
			}
		}(s)
	}
}
