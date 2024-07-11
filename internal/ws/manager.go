package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/centrifugal/centrifuge"
	"log/slog"
	"net/http"
	"time"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

type Manager struct {
	log      *slog.Logger
	node     *centrifuge.Node
	handlers map[EventType]EventHandler

	commentService CommentService
}

type CommentService interface {
	CreateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
	UpdateComment(ctx context.Context, comment domain.Comment) (domain.Comment, error)
	DeleteComment(ctx context.Context, commentID string) error
}

func NewManager(log *slog.Logger, commentService CommentService) (*Manager, error) {
	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		return nil, err
	}

	m := &Manager{
		log:      log,
		node:     node,
		handlers: make(map[EventType]EventHandler),

		commentService: commentService,
	}

	m.setupEventHandlers()

	err = m.setupNode()
	if err != nil {
		return nil, err
	}

	return m, nil
}

// setupNode configures Centrifuge Node to handle all necessary events.
func (m *Manager) setupNode() error {
	m.node.OnConnecting(func(ctx context.Context, e centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		cred, _ := centrifuge.GetCredentials(ctx)
		return centrifuge.ConnectReply{
			Data: []byte(`{}`),
			// Subscribe to a personal server-side channel.
			Subscriptions: map[string]centrifuge.SubscribeOptions{
				"#" + cred.UserID: {
					EnableRecovery: true,
					EmitPresence:   true,
					EmitJoinLeave:  true,
					PushJoinLeave:  true,
				},
			},
		}, nil
	})

	m.node.OnConnect(func(client *centrifuge.Client) {
		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 60,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			m.log.Debug("subscribe event", slog.String("channel", e.Channel), slog.String("user_id", client.UserID()))

			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EnableRecovery: true,
					EmitPresence:   true,
					EmitJoinLeave:  true,
					PushJoinLeave:  true,
					Data:           []byte(`{"msg": "welcome"}`),
				},
			}, nil)
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			m.log.Debug(
				"publish event",
				slog.String("channel", e.Channel),
				slog.String("user_id", client.UserID()),
				slog.String("data", string(e.Data)),
			)

			if !client.IsSubscribed(e.Channel) {
				cb(centrifuge.PublishReply{}, centrifuge.ErrorPermissionDenied)
				return
			}

			var msg Event
			err := json.Unmarshal(e.Data, &msg)
			if err != nil {
				cb(centrifuge.PublishReply{}, centrifuge.ErrorBadRequest)
				return
			}

			m.log.Info("received message", slog.AnyValue(msg))

			publishReply, err := m.routeEvent(clientMessage{
				Event:        msg,
				Client:       client,
				PublishEvent: e,
			})
			if err != nil {
				switch {
				case errors.Is(err, ErrEventNotSupported):
					cb(centrifuge.PublishReply{}, centrifuge.ErrorMethodNotFound)
				default:
					cb(centrifuge.PublishReply{}, centrifuge.ErrorInternal)
				}
				return
			}

			cb(publishReply, err)
		})

		client.OnPresence(func(e centrifuge.PresenceEvent, cb centrifuge.PresenceCallback) {
			if !client.IsSubscribed(e.Channel) {
				cb(centrifuge.PresenceReply{}, centrifuge.ErrorPermissionDenied)
				return
			}
			cb(centrifuge.PresenceReply{}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			m.log.Debug("unsubscribe event", slog.String("channel", e.Channel), slog.String("user_id", client.UserID()))
		})

		client.OnAlive(func() {
			m.log.Debug("alive event", slog.String("user_id", client.UserID()))
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			m.log.Debug("disconnect event", slog.String("user_id", client.UserID()))
		})
	})

	if err := m.node.Run(); err != nil {
		return err
	}

	return nil
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventCreateComment] = m.handleCreateComment
}

// routeEvent routes the event to the correct handler
//
// It will return the reply from the handler
// If the event is not supported, it will return an error
func (m *Manager) routeEvent(msg clientMessage) (centrifuge.PublishReply, error) {
	log := m.log.With(slog.String("event", string(msg.Event.Type)))

	if handler, ok := m.handlers[msg.Event.Type]; ok {
		reply, err := handler(msg)
		if err != nil {
			log.Error("error handling event", logger.Err(err))
			return centrifuge.PublishReply{}, fmt.Errorf("error handling event: %w", err)
		}
		return reply, nil
	} else {
		return centrifuge.PublishReply{}, ErrEventNotSupported
	}
}

// WebsocketHandler returns a http.Handler that can be used to upgrade HTTP
func (m *Manager) WebsocketHandler() http.Handler {
	return centrifuge.NewWebsocketHandler(m.node, centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool {
			originHeader := r.Header.Get("Origin")
			if originHeader == "" || originHeader == "null" {
				return true
			}
			return originHeader == "http://localhost:3000"
		}})
}
