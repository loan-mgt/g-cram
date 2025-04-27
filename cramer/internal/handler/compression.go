package handler

import (
	"encoding/json"
	"fmt"
	"loan-mgt/cramer/internal/service"
	"os"
)

func HandleCompression(body []byte, amqpConn *service.AMQPConnection) {

	var msg struct {
		MediaId      string `json:"mediaId"`
		UserId       string `json:"userId"`
		Timestamp    int64  `json:"timestamp"`
		CreationDate string `json:"creationDate"`
	}

	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	path := fmt.Sprintf("%s_%d/%s", msg.UserId, msg.Timestamp, msg.MediaId)

	err = os.MkdirAll(fmt.Sprintf("/tmp/out/%s_%d", msg.UserId, msg.Timestamp), 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputSize, err := service.CreateVideoWithMetadata(path, msg.CreationDate)
	if err != nil {
		fmt.Printf("Error creating video with metadata: %s\n", err)
		return
	}

	fmt.Println("Compressed video created successfully")

	err = amqpConn.SendRequest(msg.MediaId, msg.UserId, msg.Timestamp, outputSize)
	if err != nil {
		fmt.Printf("Error uploading video: %s\n", err)
		return
	}

}
