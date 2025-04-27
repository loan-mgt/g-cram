package handler

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// CreateItem creates a new item
func (h *APIHandler) GetImage(c *gin.Context) {
	var payload struct {
		BaseURL string `json:"baseUrl"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isGoogleStorageURL(payload.BaseURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BaseURL"})
		return
	}

	req, err := http.NewRequest("GET", payload.BaseURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Assume we have a way to get the user's token
	req.Header.Set("Authorization", c.Request.Header.Get("Authorization"))

	fmt.Println(req.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == http.StatusForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid or expired BaseURL"})
		return
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.Header().Set("Content-Disposition", resp.Header.Get("Content-Disposition"))
	c.Writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	c.Writer.WriteHeader(resp.StatusCode)
	c.Writer.Write(buf)
}

func (h *APIHandler) GetVideo(c *gin.Context) {
	var payload struct {
		BaseURL string `json:"baseUrl"`
		Id      string `json:"id"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isGoogleStorageURL(payload.BaseURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid BaseURL"})
		return
	}

	baseURL := payload.BaseURL + "=dv"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", c.Request.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Writer.Header().Set("Content-Disposition", resp.Header.Get("Content-Disposition"))
	c.Writer.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	c.Writer.Write(buf)
}

func isGoogleStorageURL(url string) bool {
	return regexp.MustCompile(`^https:\/\/[a-zA-Z0-9-]*\.googleusercontent\.com\/.*$`).MatchString(url)
}
