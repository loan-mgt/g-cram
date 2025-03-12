package handler

import (
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/server/ws"
	"loan-mgt/g-cram/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	amqpConn  *service.AMQPConnection
	db        *db.Store
	wsManager *ws.WebSocketManager
}

func NewAPIHandler(db *db.Store, amqpConn *service.AMQPConnection, wsManager *ws.WebSocketManager) *APIHandler {
	return &APIHandler{db: db, amqpConn: amqpConn, wsManager: wsManager}
}

// HealthCheck handles health check requests
func (h *APIHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}
