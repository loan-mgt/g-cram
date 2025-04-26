package handler

import (
	"database/sql"
	"loan-mgt/g-cram/internal/db/sqlc"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) SetUserMedia(c *gin.Context) {
	user := c.MustGet("user").(sqlc.User)

	var payload []struct {
		MediaID      string `json:"media_id"`
		CreationDate int64  `json:"creation_date"`
		Filename     string `json:"filename"`
		BaseUrl      string `json:"base_url"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(payload) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload is empty"})
		return
	}
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	for _, p := range payload {
		arg := sqlc.CreateMediaParams{
			UserID:       user.ID,
			Timestamp:    timestamp,
			MediaID:      p.MediaID,
			CreationDate: p.CreationDate,
			Filename:     p.Filename,
			BaseUrl:      p.BaseUrl,
			OldSize:      sql.NullInt64{Int64: 0, Valid: true},
			NewSize:      sql.NullInt64{Int64: 0, Valid: true},
			Done:         0,
		}

		if err := h.db.CreateMedia(c.Request.Context(), arg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// get nb of media for user
	nbMedia, err := h.db.CountUserMedia(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "nb_media": nbMedia})
}
