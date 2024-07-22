package app

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app/grpcapp"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app/httpapp"
	userclient "github.com/ARUMANDESU/uniclubs-comments-service/internal/client/user"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/grpc/commentgrpc"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/handlers"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/userservice"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/storage/mongodb"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/ws"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"log/slog"
	"sync"
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
func New(ctx context.Context, cfg config.Config, log *slog.Logger) *App {
	const op = "app.new"
	l := log.With(slog.String("op", op))

	starters := make([]Starter, 0)
	stoppers := make([]Stopper, 0)

	mongoStorage, err := mongodb.NewStorage(ctx, cfg.MongoDB)
	if err != nil {
		l.Error("failed to create mongodb storage", logger.Err(err))
		panic(err)
	}
	stoppers = append(stoppers, &mongoStorage)

	// user microservice grpc client
	userClient, err := userclient.New(log, cfg.Clients.User.Address, cfg.Clients.User.Timeout, cfg.Clients.User.RetriesCount)
	if err != nil {
		log.Error("user service client init error", logger.Err(err))
		panic(err)
	}

	userService := userservice.New(log, &mongoStorage, userClient)

	commentService := commentservice.New(commentservice.Config{
		Logger:       log,
		Provider:     &mongoStorage,
		Creator:      &mongoStorage,
		Updater:      &mongoStorage,
		Deleter:      &mongoStorage,
		UserProvider: &userService,
	})

	wsManager, err := ws.NewManager(log, commentService)
	if err != nil {
		l.Error("failed to create websocket manager", logger.Err(err))
		panic(err)
	}
	stoppers = append(stoppers, wsManager)

	handler := handlers.NewHandler(log, wsManager)
	handler.RegisterRoutes()

	httpServer := httpapp.New(cfg, log, handler.Mux)
	starters = append(starters, httpServer)
	stoppers = append(stoppers, httpServer)

	grpcServer := commentgrpc.NewServer(commentService)

	grpcApp := grpcapp.New(log, cfg.GRPC.Port, grpcServer)
	starters = append(starters, grpcApp)
	stoppers = append(stoppers, grpcApp)

	return &App{
		log:      log,
		cfg:      cfg,
		stoppers: stoppers,
		starters: starters,
	}
}

func (a *App) Start() error {
	const op = "app.start"
	log := a.log.With(slog.String("op", op))

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to signal an error
	errCh := make(chan error, 1)

	for _, s := range a.starters {
		wg.Add(1)
		go func(s Starter) {
			defer wg.Done()
			startCtx, startCancel := context.WithTimeout(ctx, a.cfg.StartTimeout)
			defer startCancel()

			if err := s.Start(startCtx); err != nil {
				log.Error("failed to start service", logger.Err(err))
				select {
				case errCh <- err:
				default:
				}
			}
		}(s)
	}

	// Wait for either all services to start or an error to occur
	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		cancel()
		return err
	case <-ctx.Done():
		return nil
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

	var wg sync.WaitGroup

	for _, s := range a.stoppers {
		wg.Add(1)
		go func(s Stopper) {
			defer wg.Done()
			if err := s.Stop(shutdownCtx); err != nil {
				log.Error("failed to stop service", logger.Err(err))
			}
		}(s)
	}

	// Wait for all services to stop
	wg.Wait()
}
