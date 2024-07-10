package ws

import (
	"encoding/json"
	"github.com/centrifugal/centrifuge"
)

// clientMessage is a struct that holds the message and the client that sent it
type clientMessage struct {
	Event        Event
	Client       *centrifuge.Client
	PublishEvent centrifuge.PublishEvent
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered depending on the type
type EventHandler func(message clientMessage) (centrifuge.PublishReply, error)

type EventType string

// Client events
const (
	// EventCreateComment is event name for new comment request
	EventCreateComment EventType = "create_comment"
)

// Server events
const (
	// EventNewComment is event name for new comment creation
	EventNewComment EventType = "new_comment"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	Type      EventType       `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
}
