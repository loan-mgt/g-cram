package handler

import (
	"context"
	"encoding/json"
	"loan-mgt/g-cram/internal/db/sqlc"
	"loan-mgt/g-cram/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) Start(c *gin.Context) {
	user := c.MustGet("user").(sqlc.User)

	accessToken, err := service.GetAccessToken(context.Background(), h.cfg, h.db, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mediaItems, err := h.db.GetMedias(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, media := range mediaItems {
		message := struct {
			Token     string `json:"token"`
			MediaId   string `json:"mediaId"`
			UserId    string `json:"userId"`
			Timestamp int64  `json:"timestamp"`
			BaseUrl   string `json:"baseUrl"`
		}{
			Token:     accessToken,
			MediaId:   media.MediaID,
			UserId:    media.UserID,
			Timestamp: media.Timestamp,
			BaseUrl:   media.BaseUrl,
		}

		messageData, err := json.Marshal(message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := h.amqpConn.SendRequest(messageData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ok"})
}
