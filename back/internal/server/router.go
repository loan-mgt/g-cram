package server

import (
	"loan-mgt/g-cram/internal/server/handler"

	"github.com/gin-gonic/gin"
)

// NewRouter sets up and configures all API routes
func NewRouter() *gin.Engine {
	router := gin.Default()

	// Create handlers
	apiHandler := handler.NewAPIHandler()

	// Define routes
	router.GET("/health", apiHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/items", apiHandler.GetItems)
		v1.GET("/items/:id", apiHandler.GetItem)
		v1.POST("/items", apiHandler.CreateItem)
		v1.PUT("/items/:id", apiHandler.UpdateItem)
		v1.DELETE("/items/:id", apiHandler.DeleteItem)
	}

	return router
}
