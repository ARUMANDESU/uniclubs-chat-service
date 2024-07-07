package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var env string
	flag.StringVar(&env, "env", ".env", "environment variables file")
	flag.Parse()

	err := godotenv.Load(env)
	if err != nil {
		log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Error(fmt.Sprintf("error loading .env file: %v", err))
	}

	cfg := config.MustLoad()
	log := logger.Setup(cfg.Env)

	log.Info("starting application",
		slog.String("env", cfg.Env),
		slog.String("port", cfg.HTTP.Address),
	)

	application := app.New(*cfg, log)

	go func() {
		if err := application.HTTPSvr.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error: %v", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		}
		log.Info("stopped serving new connections")

	}()

	log.Info("application started")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	defer log.Info("application stopped", slog.String("signal", sign.String()))
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.Stop()
}
