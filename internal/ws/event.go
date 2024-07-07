package ws

import "encoding/json"

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered depending on the type
type EventHandler func(event Event, c *Client) error

// Client events
const (
	// EventCreateComment is event name for new comment request
	EventCreateComment = "create_comment"
)

// Server events
const (
	// EventNewComment is event name for new comment creation
	EventNewComment = "new_comment"
)
