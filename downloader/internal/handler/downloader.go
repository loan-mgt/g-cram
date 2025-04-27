package handler

import (
	"encoding/json"
	"fmt"
	"loan-mgt/uploader/internal/service"
)

func HandleDownload(body []byte, conn *service.AMQPConnection) {

	var msg service.Msg

	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileSize, err := service.DownloadVideo(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Video downloaded successfully")

	fmt.Println("media_id", msg.MediaId, "file_size", fileSize, "timestamp", msg.Timestamp, "user_id", msg.UserId)

	err = conn.SendNotificationRequest(msg.MediaId, msg.UserId, msg.Timestamp, fileSize)
	if err != nil {
		fmt.Printf("Error uploading video: %s\n", err)
		return
	}

}
