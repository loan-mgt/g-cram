package handler

import (
	"io"
	"loan-mgt/g-cram/internal/db/sqlc"
	"loan-mgt/g-cram/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) AddSubscriptionToUser(c *gin.Context) {
	user := c.MustGet("user").(sqlc.User)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := service.AddSubscriptionToUser(c.Request.Context(), h.db, user.ID, string(body)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *APIHandler) RemoveSubscriptionFromUser(c *gin.Context) {
	user := c.MustGet("user").(sqlc.User)

	if err := service.RemoveSubscriptionFromUser(c.Request.Context(), h.db, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{})
}
