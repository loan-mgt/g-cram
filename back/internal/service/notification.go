package service

import (
	"encoding/json"
	"fmt"
	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/server/ws"
)

func StartUploadDoneListener(ws *ws.WebSocketManager, db *db.Store, cfg *config.Config) {

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
			UserId   string `json:"userId"`
		}

		if err := json.Unmarshal(d.Body, &payload); err != nil {
			fmt.Println("Error unmarshal ling message:", err)
			continue
		}
		if err := ws.SendMessageToClient(payload.UserId, "Success"); err != nil {
			fmt.Println("Error sending message to WebSocket:", err)
		}

		err = PushNotification(db, cfg, payload.UserId, fmt.Sprintf("Notification: %s", payload.Filename), "job-done:"+payload.Filename)
		if err != nil {
			fmt.Println("Error sending push notification:", err)
		}

	}

}
