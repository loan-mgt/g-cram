package handler

import (
	"context"
	"database/sql"
	"loan-mgt/g-cram/internal/db/sqlc"
	"loan-mgt/g-cram/internal/server/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow requests from specific domains
		allowedOrigins := []string{"http://localhost:8080"}
		for _, origin := range allowedOrigins {
			if r.Header.Get("Origin") == origin {
				return true
			}
		}
		return false
	},
}

// WebSocketHandler handles WebSocket connections.
func (h *APIHandler) WebSocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	// Extract client ID and token from query parameters
	token := c.Request.URL.Query().Get("token")
	id := c.Request.URL.Query().Get("id")

	// Check if the user exists in the database
	_, err = h.db.GetUser(context.Background(), id)
	if err != nil {
		// Create user if not exists
		arg := sqlc.CreateUserParams{
			ID:        id,
			Token:     sql.NullString{String: token, Valid: true},
			Websocket: sql.NullString{String: conn.RemoteAddr().String(), Valid: true},
		}
		if err = h.db.CreateUser(context.Background(), arg); err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			conn.Close()
			return
		}
	} else {
		// Update user token and WebSocket connection
		arg := sqlc.UpdateUserTokenAndWebsocketParams{
			ID:        id,
			Token:     sql.NullString{String: token, Valid: true},
			Websocket: sql.NullString{String: conn.RemoteAddr().String(), Valid: true},
		}
		if err = h.db.UpdateUserTokenAndWebsocket(context.Background(), arg); err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			conn.Close()
			return
		}
	}

	// Register WebSocket connection with the manager
	h.wsManager.AddClient(id, conn)

	// Send welcome message
	if err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from cramer")); err != nil {
		log.Println("WebSocket write error:", err)
		h.wsManager.RemoveClient(id)
		conn.Close()
	}
}
