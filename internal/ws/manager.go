package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/centrifugal/centrifuge"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

const exampleChannel = "chat:index"

// Check whether channel is allowed for subscribing. In real case permission
// check will probably be more complex than in this example.
func channelSubscribeAllowed(channel string) bool {
	return channel == exampleChannel
}

type Manager struct {
	sync.RWMutex

	log      *slog.Logger
	node     *centrifuge.Node
	handlers map[EventType]EventHandler
}

func NewManager(log *slog.Logger) (*Manager, error) {
	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		return nil, err
	}

	m := &Manager{
		log:      log,
		node:     node,
		handlers: make(map[EventType]EventHandler),
	}

	m.setupEventHandlers()
	m.setupNode()

	return m, nil
}

// setupNode configures Centrifuge Node to handle all necessary events.
func (m *Manager) setupNode() {
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
		transport := client.Transport()
		log.Printf("[user %s] connected via %s with protocol: %s", client.UserID(), transport.Name(), transport.Protocol())

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			log.Printf("[user %s] connection is going to expire, refreshing", client.UserID())

			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 60,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			log.Printf("[user %s] subscribes on %s", client.UserID(), e.Channel)

			if !channelSubscribeAllowed(e.Channel) {
				cb(centrifuge.SubscribeReply{}, centrifuge.ErrorPermissionDenied)
				return
			}

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
			log.Printf("[user %s] publishes into channel %s: %s", client.UserID(), e.Channel, string(e.Data))

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

		client.OnMessage(func(e centrifuge.MessageEvent) {
			log.Printf("[user %s] message received: %s", client.UserID(), string(e.Data))
		})

		client.OnRPC(func(e centrifuge.RPCEvent, cb centrifuge.RPCCallback) {
			log.Printf("[user %s] sent RPC, data: %s, method: %s", client.UserID(), string(e.Data), e.Method)
			switch e.Method {
			case "getCurrentYear":
				cb(centrifuge.RPCReply{Data: []byte(`{"year": "2020"}`)}, nil)
			default:
				cb(centrifuge.RPCReply{}, centrifuge.ErrorMethodNotFound)
			}
		})

		client.OnPresence(func(e centrifuge.PresenceEvent, cb centrifuge.PresenceCallback) {
			log.Printf("[user %s] calls presence on %s", client.UserID(), e.Channel)

			if !client.IsSubscribed(e.Channel) {
				cb(centrifuge.PresenceReply{}, centrifuge.ErrorPermissionDenied)
				return
			}
			cb(centrifuge.PresenceReply{}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			log.Printf("[user %s] unsubscribed from %s: %s", client.UserID(), e.Channel, e.Reason)
		})

		client.OnAlive(func() {
			log.Printf("[user %s] connection is still active", client.UserID())
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			log.Printf("[user %s] disconnected: %s", client.UserID(), e.Reason)
		})
	})

	if err := m.node.Run(); err != nil {
		log.Fatal(err)
	}
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventCreateComment] = func(msg clientMessage) (centrifuge.PublishReply, error) {
		msg.Event.Timestamp = time.Now().Unix()
		data, _ := json.Marshal(msg.Event)

		result, err := m.node.Publish(
			msg.PublishEvent.Channel, data,
			centrifuge.WithHistory(300, time.Minute),
			centrifuge.WithClientInfo(msg.PublishEvent.ClientInfo),
		)
		if err != nil {
			return centrifuge.PublishReply{}, fmt.Errorf("error publishing message: %w", err)
		}

		return centrifuge.PublishReply{Result: &result}, nil
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
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

func (m *Manager) WebsocketHandler() http.Handler {
	return centrifuge.NewWebsocketHandler(m.node, centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool {
			originHeader := r.Header.Get("Origin")
			log.Printf("origin header: %s", originHeader)
			if originHeader == "" || originHeader == "null" {
				return true
			}
			return originHeader == "http://localhost:3000"
		}})
}
