package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// send request on rabbit mq to comrpess
func (h *APIHandler) Start(c *gin.Context) {
	var payload []struct {
		Id           string `json:"id"`
		CreationDate string `json:"creationDate"`
		Name         string `json:"name"`
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

	token := c.Request.Header.Get("Authorization")

	for _, p := range payload {
		videoPath := fmt.Sprintf("/tmp/in/%s.mp4", p.Id)
		fmt.Println(videoPath)
		err = h.amqpConn.SendRequest(token, p.Id, p.CreationDate, p.Name, videoPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ok"})
}
