package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ARUMANDESU/uniclubs-comments-service/internal/services/commentservice"
	"github.com/centrifugal/centrifuge"
)

// handleCreateComment is an event handler that is triggered when a client sends a create_comment event
func (m *Manager) handleCreateComment(message clientMessage) (centrifuge.PublishReply, error) {
	var input struct {
		Body   string `json:"body"`
		PostID string `json:"post_id"`
	}
	err := json.Unmarshal(message.Event.Payload, &input)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	userID, err := strconv.ParseInt(message.PublishEvent.ClientInfo.UserID, 10, 64)
	if err != nil {
		return centrifuge.PublishReply{}, fmt.Errorf("error converting UserID to int64: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdComment, err := m.commentService.Create(ctx, commentservice.CreateCommentDTO{
		Body:   input.Body,
		PostID: input.PostID,
		UserID: userID,
	})
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
	var input struct {
		CommentID string `json:"comment_id"`
	}
	err := json.Unmarshal(message.Event.Payload, &input)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	userID, err := strconv.ParseInt(message.PublishEvent.ClientInfo.UserID, 10, 64)
	if err != nil {
		return centrifuge.PublishReply{}, fmt.Errorf("error converting UserID to int64: %w", err)
	}

	err = m.commentService.Delete(context.TODO(), commentservice.DeleteCommentDTO{
		CommentID: input.CommentID,
		UserID:    userID,
	})
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
	var input struct {
		CommentID string `json:"comment_id"`
		Body      string `json:"body"`
	}
	err := json.Unmarshal(message.Event.Payload, &input)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	userID, err := strconv.ParseInt(message.PublishEvent.ClientInfo.UserID, 10, 64)
	if err != nil {
		return centrifuge.PublishReply{}, fmt.Errorf("error converting UserID to int64: %w", err)
	}

	updatedComment, err := m.commentService.Update(context.TODO(), commentservice.UpdateCommentDTO{
		CommentID: input.CommentID,
		Body:      input.Body,
		UserID:    userID,
	})
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	payload, err := json.Marshal(updatedComment)
	if err != nil {
		return centrifuge.PublishReply{}, err
	}

	event := Event{
		Type:      EventEditComment,
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
