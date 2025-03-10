package server

import (
	"loan-mgt/g-cram/internal/server/handler"
	"loan-mgt/g-cram/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter sets up and configures all API routes
func NewRouter(amqpConn *service.AMQPConnection) *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	// Create handlers
	apiHandler := handler.NewAPIHandler(amqpConn)

	// Define routes
	router.GET("/health", apiHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.POST("/get-image", apiHandler.GetImage)
		v1.POST("/get-video", apiHandler.GetVideo)
		v1.POST("/start", apiHandler.Start)
	}

	return router
}
