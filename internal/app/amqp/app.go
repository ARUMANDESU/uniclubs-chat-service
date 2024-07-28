package amqpapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/rabbitmq"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
)

type App struct {
	log        *slog.Logger
	amqp       Amqp
	usrService UserService
}

//go:generate mockery --name Amqp
type Amqp interface {
	Consume(queue string, routingKey string, handler func(msg amqp091.Delivery) error) error
	Close() error
}

//go:generate mockery --name UserService
type UserService interface {
	Update(ctx context.Context, user domain.User) error
}

func New(log *slog.Logger, userService UserService, amqp Amqp) *App {
	return &App{
		log:        log,
		amqp:       amqp,
		usrService: userService,
	}
}

func (a *App) Start(_ context.Context, _ func(error)) {
	a.consumeMessages(rabbitmq.UserEventsQueue, rabbitmq.UserUpdatedEventRoutingKey, a.HandleUpdateUser)
}

func (a *App) consumeMessages(queue, routingKey string, handler rabbitmq.Handler) {
	go func() {
		const op = "amqp.app.consumeMessages"
		log := a.log.With(slog.String("op", op))

		err := a.amqp.Consume(queue, routingKey, handler)
		if err != nil {
			log.Error("failed to consume ", logger.Err(err))
		}
	}()
}

func (a *App) Stop(_ context.Context) error {
	const op = "amqp.app.shutdown"

	a.log.With(slog.String("op", op)).Info("shutting down amqp app")
	err := a.amqp.Close()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) HandleUpdateUser(msg amqp091.Delivery) error {
	const op = "amqp.app.handle-update-user"
	log := a.log.With(slog.String("op", op))

	var user domain.User

	err := json.Unmarshal(msg.Body, &user)
	if err != nil {
		log.Error("failed to unmarshal message", logger.Err(err))
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = a.usrService.Update(ctx, user)
	if err != nil {
		log.Error("failed to update user", logger.Err(err))
		return err
	}

	return nil
}
