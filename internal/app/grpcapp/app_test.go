package grpcapp

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/grpc/commentgrpc"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestAppStart_Success(t *testing.T) {
	grpcServer := commentgrpc.NewServer(nil)
	app := New(logger.Plug(), 0, &grpcServer)

	go app.Start(context.Background(), func(err error) {
		assert.NoError(t, err)
	})

	assert.NotEmpty(t, app.gRPCServer.GetServiceInfo())

	err := app.Stop(context.Background())
	assert.NoError(t, err)
}

func TestAppStart_Failure_PortAlreadyInUse(t *testing.T) {
	lsn, _ := net.Listen("tcp", ":0")

	app := New(logger.Plug(), lsn.Addr().(*net.TCPAddr).Port, nil)

	errCh := make(chan error)
	go app.Start(context.Background(), func(err error) {
		errCh <- err
	})

	err := <-errCh
	assert.Error(t, err)
}
