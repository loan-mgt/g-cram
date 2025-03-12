package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection.
type Client struct {
	ID     string
	conn   *websocket.Conn
	manager *WebSocketManager
}

// NewClient creates a new WebSocket client.
func NewClient(id string, conn *websocket.Conn, manager *WebSocketManager) *Client {
	return &Client{
		ID:      id,
		conn:    conn,
		manager: manager,
	}
}

// Listen listens for messages from the client.
func (c *Client) Listen() {
	defer func() {
		c.manager.RemoveClient(c.ID)
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		log.Printf("Received message from client %s: %s", c.ID, msg)
	}
}
