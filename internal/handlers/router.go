package handlers

func (h *Handler) RegisterRoutes() {
	h.Mux.Handle("/connection/websocket", authMiddleware(h.wsHandler.WebsocketHandler()))
}
