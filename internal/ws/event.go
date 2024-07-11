package ws

import (
	"encoding/json"
	"github.com/centrifugal/centrifuge"
)

// clientMessage is a struct that contains the event and the client that triggered the event
type clientMessage struct {
	Event        Event
	Client       *centrifuge.Client
	PublishEvent centrifuge.PublishEvent
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered depending on the type
type EventHandler func(message clientMessage) (centrifuge.PublishReply, error)

type EventType string

// Client events which receive from the client
const (
	EventCreateComment EventType = "create_comment"
	EventUpdateComment EventType = "update_comment"
	EventDeleteComment EventType = "delete_comment"
)

// Server events which are sent to the client
const (
	EventNewComment    EventType = "new_comment"
	EventEditComment   EventType = "edit_comment"
	EventRemoveComment EventType = "remove_comment"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	Type      EventType       `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
}
