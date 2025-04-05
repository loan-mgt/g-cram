package handler

import (
	"loan-mgt/g-cram/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

func (h *APIHandler) SetUserMedia(c *gin.Context) {
	user := c.MustGet("user").(sqlc.User)

	// clear user media
	h.db.ClearUserTmpMedia(c.Request.Context(), user.ID)

}
