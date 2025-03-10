package handler

import (
	"encoding/json"
	"fmt"
	"loan-mgt/uploader/internal/service"
)

func HandleUpload(body []byte) {
	var msg struct {
		Token     string `json:"token"`
		VideoPath string `json:"videoPath"`
		FileName  string `json:"fileName"`
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

}
