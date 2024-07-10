package handlers

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	log       *slog.Logger
	wsHandler WebsocketHandler

	Mux *http.ServeMux
}

type WebsocketHandler interface {
	WebsocketHandler() http.Handler
}

func NewHandler(log *slog.Logger, wsHandler WebsocketHandler) *Handler {
	return &Handler{
		log:       log,
		Mux:       http.NewServeMux(),
		wsHandler: wsHandler,
	}
}
