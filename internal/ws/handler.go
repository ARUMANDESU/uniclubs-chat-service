package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-comments-service/internal/domain"
	"github.com/centrifugal/centrifuge"
	"time"
)

// handleCreateComment is an event handler that is triggered when a client sends a create_comment event
func (m *Manager) handleCreateComment(message clientMessage) (centrifuge.PublishReply, error) {
	var comment domain.Comment
	err := json.Unmarshal(message.Event.Payload, &comment)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	// Create the comment
	createdComment, err := m.commentService.CreateComment(context.TODO(), comment)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	// Marshal the created comment into a json.RawMessage
	payload, err := json.Marshal(createdComment)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	// Create a new event with the created comment as the payload
	event := Event{
		Type:      EventNewComment,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(event)

	result, err := m.node.Publish(
		message.PublishEvent.Channel, data,
		centrifuge.WithHistory(300, time.Minute),
		centrifuge.WithClientInfo(message.PublishEvent.ClientInfo),
	)
	if err != nil {
		return centrifuge.PublishReply{}, fmt.Errorf("error publishing message: %w", err)
	}

	return centrifuge.PublishReply{Result: &result}, nil
}
