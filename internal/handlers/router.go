package handlers

func (h *Handler) RegisterRoutes() {
	h.Mux.Handle("/connection/websocket", h.wsHandler.WebsocketHandler())
}
