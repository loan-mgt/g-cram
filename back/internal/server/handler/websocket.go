package handler

import (
	"context"
	"database/sql"
	"loan-mgt/g-cram/internal/db/sqlc"
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

func (h *APIHandler) WebSocket(c *gin.Context) {
	w, r := c.Writer, c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()
	//?token=${accessToken}&id=${urlFragment.split("client_id=")[1].split("&")[0]}
	// extract info create user if not existe, else update
	token := c.Request.URL.Query().Get("token")
	id := c.Request.URL.Query().Get("id")

	_, err = h.db.GetUser(context.Background(), id)
	if err != nil {
		// create user
		arg := sqlc.CreateUserParams{
			ID:        id,
			Token:     sql.NullString{String: token, Valid: true},
			Websocket: sql.NullString{String: conn.RemoteAddr().String(), Valid: true},
		}
		err = h.db.CreateUser(context.Background(), arg)
		if err != nil {
			err = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			if err != nil {
				log.Println("write:", err)
				return
			}
			return
		}
	} else {
		// update user
		arg := sqlc.UpdateUserTokenAndWebsocketParams{
			ID:        id,
			Token:     sql.NullString{String: token, Valid: true},
			Websocket: sql.NullString{String: conn.RemoteAddr().String(), Valid: true},
		}
		err = h.db.UpdateUserTokenAndWebsocket(context.Background(), arg)
		if err != nil {
			err = conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			if err != nil {
				log.Println("write:", err)
				return
			}
			return
		}
	}

	// send message
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from cramer"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	// receive message
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}
	log.Printf("recv: %s", message)
}
