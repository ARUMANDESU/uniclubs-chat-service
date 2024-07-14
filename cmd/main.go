package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/app"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/joho/godotenv"
	"log/slog"
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

	application := app.New(context.Background(), *cfg, log)

	err = application.Start()
	if err != nil {
		log.Error("failed to start application", logger.Err(err))
		os.Exit(1)
	}

	log.Info("application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))
	defer log.Info("application stopped", slog.String("signal", sign.String()))

	application.Stop()
}
