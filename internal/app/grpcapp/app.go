package grpcapp

import (
	"context"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/grpc/commentgrpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, server commentgrpc.Server) *App {
	gRPCServer := grpc.NewServer()

	commentgrpc.Register(gRPCServer, server)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Start(context.TODO()); err != nil {
		panic(err)
	}
}

func (a *App) Start(_ context.Context) error {
	const op = "app.grpc.run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("grpc port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		log.Info("gRPCServer is running", slog.String("addr", l.Addr().String()))
		if err := a.gRPCServer.Serve(l); err != nil {
			log.Error("gRPCServer error", slog.String("err", err.Error()))
		}
	}()

	return nil
}

func (a *App) Stop(_ context.Context) error {
	const op = "app.grpc.stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC Server")
	a.gRPCServer.GracefulStop()

	return nil
}
