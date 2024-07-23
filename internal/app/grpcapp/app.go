package grpcapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/grpc/commentgrpc"
	commentv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/comment"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, server commentv1.CommentServer) *App {
	gRPCServer := grpc.NewServer()

	commentgrpc.Register(gRPCServer, server)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// Start starts the gRPC server
//
// Panics if an error occurs while starting the gRPC server
func (a *App) Start(_ context.Context, callback func(error)) {
	const op = "app.grpc.start"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		callback(err)
		return
	}

	go func() {
		log.Info("gRPC server is running", slog.String("addr", l.Addr().String()))
		if err := a.gRPCServer.Serve(l); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error("gRPC server error", slog.String("err", err.Error()))
			callback(err)
		}
	}()

}

func (a *App) Stop(_ context.Context) error {
	const op = "app.grpc.stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC Server")
	a.gRPCServer.GracefulStop()

	return nil
}
