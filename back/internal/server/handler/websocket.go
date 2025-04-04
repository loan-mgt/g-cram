package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections.
func (h *APIHandler) WebSocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Extract client ID and token from query parameters
	id := c.Request.URL.Query().Get("id")

	// Register WebSocket connection with the manager
	h.wsManager.AddClient(id, conn)

	// Send welcome message
	if err = conn.WriteMessage(websocket.TextMessage, []byte("{\"msg\": \"Hello from cramer\"}")); err != nil {
		log.Println("WebSocket write error:", err)
		h.wsManager.RemoveClient(id)
		conn.Close()
	}
}
