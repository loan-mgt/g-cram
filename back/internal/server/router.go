package server

import (
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/server/handler"
	"loan-mgt/g-cram/internal/server/middleware"
	"loan-mgt/g-cram/internal/server/ws"
	"loan-mgt/g-cram/internal/service"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// NewRouter sets up and configures all API routes
func NewRouter(store *db.Store, amqpConn *service.AMQPConnection, ws *ws.WebSocketManager, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.FrontUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow requests from specific domains
			allowedOrigins := []string{cfg.FrontUrl}
			for _, origin := range allowedOrigins {
				if r.Header.Get("Origin") == origin {
					return true
				}
			}
			return false
		},
	}

	// Create handlers
	apiHandler := handler.NewAPIHandler(store, amqpConn, ws, cfg, &upgrader)

	// Create middleware context
	mc := middleware.NewMiddleWareContextr(store, cfg)

	// Define routes
	router.GET("/health", apiHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/get-image", apiHandler.GetImage)
		v1.POST("/get-video", apiHandler.GetVideo)
		v1.POST("/start", mc.AuthMiddleware(), apiHandler.Start)
		v1.GET("/user", mc.AuthMiddleware(), apiHandler.GetUser)
		v1.GET("/job", mc.AuthMiddleware(), apiHandler.GetJob)
		v1.POST("/user", apiHandler.InitUser)
		v1.PATCH(("/media"), mc.AuthMiddleware(), apiHandler.SetUserMedia)
		v1.POST("/user/subscription", mc.AuthMiddleware(), apiHandler.AddSubscriptionToUser)
		v1.DELETE("/user/subscription", mc.AuthMiddleware(), apiHandler.RemoveSubscriptionFromUser)

		v1.GET("/ws", apiHandler.WebSocket)
	}

	return router
}
