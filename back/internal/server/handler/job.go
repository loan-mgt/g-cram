package handler

import (
	"context"
	"loan-mgt/g-cram/internal/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) GetJob(c *gin.Context) {
	dbUser := c.MustGet("user").(sqlc.User)

	// Get tokens using authorization code
	jobs, err := h.db.GetUserJob(context.Background(), dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	type job struct {
		Timestamp   int64   `json:"timestamp"`
		NbMedia     int64   `json:"nb_media"`
		NbMediaDone int64   `json:"nb_media_done"`
		OldSize     float64 `json:"old_size"`
		NewSize     float64 `json:"new_size"`
	}

	var jobsList []job
	for _, j := range jobs {
		jobsList = append(jobsList, job{
			Timestamp:   j.Timestamp,
			NbMedia:     j.NbMedia,
			NbMediaDone: j.NbMediaDone,
			OldSize:     j.OldSize.Float64,
			NewSize:     j.NewSize.Float64,
		})
	}

	c.JSON(http.StatusOK, jobsList)
}
