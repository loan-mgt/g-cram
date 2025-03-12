package websocket

// Hub manages multiple WebSocket clients.
type Hub struct {
	manager *WebSocketManager
}

// NewHub creates a new WebSocket hub.
func NewHub(manager *WebSocketManager) *Hub {
	return &Hub{manager: manager}
}

// Broadcast sends a message to all clients.
func (h *Hub) Broadcast(message string) {
	h.manager.BroadcastMessage(message)
}
