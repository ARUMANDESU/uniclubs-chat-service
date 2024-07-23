package httpapp

import (
	"context"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/thejerf/slogassert"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestServerStartAndStop(t *testing.T) {
	cfg := config.Config{
		HTTP: config.HTTP{
			Address:     ":9090",
			Timeout:     10 * time.Second,
			IdleTimeout: 10 * time.Second,
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := New(cfg, logger.Plug(), handler)

	go func() {
		server.Start(context.Background(), func(err error) {
			assert.NoError(t, err)
		})
	}()

	resp, err := http.Get("http://localhost:9090")
	if err != nil {
		t.Errorf("Failed to send request to server: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	if err := server.Stop(context.Background()); err != nil {
		t.Errorf("Failed to stop server: %v", err)
	}
}

func TestServerStartFailure(t *testing.T) {
	cfg := config.Config{
		HTTP: config.HTTP{
			Address: "invalid_address",
		},
	}

	handler := slogassert.New(t, slog.LevelWarn, nil)
	log := slog.New(handler)

	server := New(cfg, log, nil)

	done := make(chan struct{})
	go func() {
		server.Start(context.Background(), func(err error) {
			assert.Error(t, err)
			close(done)
		})
	}()
	<-done
}
