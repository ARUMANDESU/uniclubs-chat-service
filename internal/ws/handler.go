package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice"
	"github.com/centrifugal/centrifuge"
)

// handleCreateComment is an event handler that is triggered when a client sends a create_comment event
func (m *Manager) handleCreateComment(message clientMessage) (centrifuge.PublishReply, error) {
	var dto commentservice.CreateCommentDTO
	err := json.Unmarshal(message.Event.Payload, &dto)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	// Create the comment
	createdComment, err := m.commentService.Create(context.TODO(), dto)
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

func (m *Manager) handleDeleteComment(message clientMessage) (centrifuge.PublishReply, error) {
	var dto commentservice.DeleteCommentDTO
	err := json.Unmarshal(message.Event.Payload, &dto)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	err = m.commentService.Delete(context.TODO(), dto)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	event := Event{
		Type:      EventRemoveComment,
		Payload:   message.Event.Payload,
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

func (m *Manager) handleUpdateComment(message clientMessage) (centrifuge.PublishReply, error) {
	var dto commentservice.UpdateCommentDTO
	err := json.Unmarshal(message.Event.Payload, &dto)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	updatedComment, err := m.commentService.Update(context.TODO(), dto)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	payload, err := json.Marshal(updatedComment)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	event := Event{
		Type:      EventUpdateComment,
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
