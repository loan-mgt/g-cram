package service

import (
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/server/ws"
)

func StartListener(ws *ws.WebSocketManager) {

	notification, err := NewAMQPConnection(config.New(), "notification")
	if err != nil {
		panic(err)
	}

	defer notification.Conn.Close()

	// Consume messages from queue
	msgs, err := notification.Consume()
	if err != nil {
		fmt.Println(err)
		return
	}

	for d := range msgs {
		var payload struct {
			Token    string `json:"token"`
			Filename string `json:"filename"`
		}
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			fmt.Println("Error unmarshal ling message:", err)
			continue
		}
		if err := ws.SendMessageToClient(payload.Token, "Success"); err != nil {
			fmt.Println("Error sending message to WebSocket:", err)
		}

	}

}
