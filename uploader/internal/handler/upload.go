package handler

import (
	"encoding/json"
	"fmt"
	"loan-mgt/uploader/internal/service"
)

func HandleUpload(body []byte, conn *service.AMQPConnection) {
	var msg struct {
		MediaId   string `json:"mediaId"`
		UserId    string `json:"userId"`
		Timestamp int64  `json:"timestamp"`
		Token     string `json:"token"`
		FileName  string `json:"fileName"`
	}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	videoPath := fmt.Sprintf("/tmp/out/%s_%d/%s.mp4", msg.UserId, msg.Timestamp, msg.MediaId)

	err = service.UploadVideo(msg.Token, videoPath, msg.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Video uploaded successfully")

	err = conn.SendNotificationRequest(msg.MediaId, msg.UserId, msg.Timestamp)
	if err != nil {
		fmt.Printf("Error uploading video: %s\n", err)
		return
	}

}
