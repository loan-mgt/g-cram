package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mc *MiddleWareContext) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenHash, err := c.Cookie("th")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		dbUser, err := mc.db.GetUserByTokenHash(c.Request.Context(), sql.NullString{
			String: tokenHash,
			Valid:  true,
		})
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.Set("user", dbUser)

		c.Next()
	}
}
