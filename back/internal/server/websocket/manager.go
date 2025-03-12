package websocket

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketManager manages active WebSocket connections.
type WebSocketManager struct {
	clients map[string]*Client
	mu      sync.Mutex
}

// NewWebSocketManager creates a new instance of WebSocketManager.
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[string]*Client),
	}
}

// AddClient registers a new WebSocket connection.
func (wm *WebSocketManager) AddClient(id string, conn *websocket.Conn) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	client := NewClient(id, conn, wm)
	wm.clients[id] = client
	go client.Listen()
}

// RemoveClient removes a WebSocket connection.
func (wm *WebSocketManager) RemoveClient(id string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if client, exists := wm.clients[id]; exists {
		client.conn.Close()
		delete(wm.clients, id)
	}
}

// SendMessageToClient sends a message to a specific WebSocket client.
func (wm *WebSocketManager) SendMessageToClient(id, message string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	client, exists := wm.clients[id]
	if !exists {
		return fmt.Errorf("client not found")
	}

	return client.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// BroadcastMessage sends a message to all connected clients.
func (wm *WebSocketManager) BroadcastMessage(message string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for _, client := range wm.clients {
		_ = client.conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
