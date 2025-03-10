package handler

import (
	"encoding/json"
	"fmt"
	"loan-mgt/cramer/internal/service"
)

func HandleCompression(body []byte, amqpConn *service.AMQPConnection) {

	var msg struct {
		Token        string `json:"token"`
		Id           string `json:"id"`
		CreationDate string `json:"creationDate"`
		Name         string `json:"name"`
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.CreateVideoWithMetadata(msg.Id, msg.CreationDate)
	if err != nil {
		fmt.Printf("Error creating video with metadata: %s\n", err)
		return
	}

	fmt.Println("Compressed video created successfully")

	err = amqpConn.SendRequest(msg.Token, msg.Name, fmt.Sprintf("/tmp/out/%s.mp4", msg.Id))
	if err != nil {
		fmt.Printf("Error uploading video: %s\n", err)
		return
	}

}
