package handler

import (
	"context"
	"database/sql"
	"loan-mgt/g-cram/internal/db/sqlc"
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
	token := c.Request.URL.Query().Get("token")
	id := c.Request.URL.Query().Get("id")

	// Check if the user exists in the database
	_, err = h.db.GetUser(context.Background(), id)
	if err != nil {
		arg := sqlc.CreateUserParams{
			ID:    id,
			Token: sql.NullString{String: token, Valid: true},
		}
		if err = h.db.CreateUser(context.Background(), arg); err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			conn.Close()
			return
		}
	} else {
		arg := sqlc.UpdateUserTokenParams{
			ID:    id,
			Token: sql.NullString{String: token, Valid: true},
		}
		if err = h.db.UpdateUserToken(context.Background(), arg); err != nil {
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
