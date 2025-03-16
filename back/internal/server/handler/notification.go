package handler

import (
	"io"
	"loan-mgt/g-cram/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) AddSubscriptionToUser(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := service.AddSubscriptionToUser(c.Request.Context(), h.db, userId, string(body)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *APIHandler) RemoveSubscriptionFromUser(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
	}

	if err := service.RemoveSubscriptionFromUser(c.Request.Context(), h.db, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}
