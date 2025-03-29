package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoogleClaims struct {
	ISS           string `json:"iss"`
	AZP           string `json:"azp"`
	AUD           string `json:"aud"`
	SUB           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	IAT           int64  `json:"iat"`
	EXP           int64  `json:"exp"`
}

func (h *APIHandler) InitUser(c *gin.Context) {
	var payload struct {
		Code string `json:"code"`
	}

	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := json.Unmarshal(jsonData, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tokens using authorization code
	tokens, err := service.GetTokens(h.cfg, payload.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract user info from ID token
	claims, err := extractUserInfo(tokens.IDToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sha := service.GetSha(h.cfg.Salt + tokens.RefreshToken)
	c.Header("Set-Cookie", fmt.Sprintf("th=%s; Max-Age=1814400; Path=/; SameSite=Lax", sha))

	// Save to database here
	// if exists upadte token else create user
	service.UpdateOrCreateUser(context.Background(), h.db, claims.SUB, tokens.RefreshToken, sha)

	c.JSON(http.StatusOK, gin.H{
		"userId":      claims.SUB,
		"userName":    claims.Name,
		"accessToken": tokens.AccessToken,
		"picture":     claims.Picture,
	})
}

func (h *APIHandler) GetUser(c *gin.Context) {
	tokenHash, err := c.Cookie("th")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dbUser, err := h.db.GetUserByTokenHash(c.Request.Context(), sql.NullString{
		String: tokenHash,
		Valid:  true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get tokens using authorization code
	tokens, err := service.GetTokens(h.cfg, dbUser.Token.String)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userId":      dbUser.ID,
		"accessToken": tokens,
	})
}

func extractUserInfo(token string) (*GoogleClaims, error) {
	// Split the token into its three parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode the payload part (the second part)
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error decoding payload: %v", err)
	}

	// Parse the payload into our claims struct
	var claims GoogleClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("error parsing claims: %v", err)
	}

	return &claims, nil
}
