package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/db/sqlc"
	"loan-mgt/g-cram/internal/service"
	"net/http"
	"time"

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

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	for _, media := range mediaItems {

		// update media to current timestamp
		if err := h.db.SetMediaTimestamp(c.Request.Context(), sqlc.SetMediaTimestampParams{
			MediaID:     media.MediaID,
			UserID:      user.ID,
			Timestamp:   timestamp,
			Timestamp_2: media.Timestamp,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
			Timestamp: timestamp,
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

		params := sqlc.SetMediaStepParams{
			MediaID:   media.MediaID,
			UserID:    user.ID,
			Timestamp: timestamp,
			Step:      media.Step + 1,
		}

		// move media to next step
		if err := h.db.SetMediaStep(c.Request.Context(), params); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	err = service.PushNotification(h.db, h.cfg, user.ID, fmt.Sprintf("Job start: for %d videos", len(mediaItems)), fmt.Sprintf("%d", timestamp))
	if err != nil {
		fmt.Println("Error sending push notification:", err)
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ok"})
}
