package handler

import (
	"encoding/json"
	"fmt"
	"loan-mgt/uploader/internal/service"
)

func HandleUpload(body []byte, conn *service.AMQPConnection) {
	var msg struct {
		Token     string `json:"token"`
		VideoPath string `json:"videoPath"`
		FileName  string `json:"fileName"`
		UserId    string `json:"userId"`
	}
	err := json.Unmarshal(body, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = service.UploadVideo(msg.Token, msg.VideoPath, msg.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Video uploaded successfully")

	err = conn.SendNotificationRequest(msg.Token, msg.FileName, msg.UserId)
	if err != nil {
		fmt.Printf("Error uploading video: %s\n", err)
		return
	}

}
