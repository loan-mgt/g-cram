package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// NewMediaItem represents a new media item to be created
type NewMediaItem struct {
	Description     string          `json:"description"`
	SimpleMediaItem SimpleMediaItem `json:"simpleMediaItem"`
}

type SimpleMediaItem struct {
	FileName    string `json:"fileName"`
	UploadToken string `json:"uploadToken"`
}
type AlbumPosition struct {
	Position            string `json:"position"`
	RelativeMediaItemId string `json:"relativeMediaItemId"`
}

func UploadVideo(token string, videoPath string, fileName string) error {
	file, err := os.Open(videoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("POST", "https://photoslibrary.googleapis.com/v1/uploads", file)
	if err != nil {
		return err
	}
	token = fmt.Sprintf("Bearer %s", token)

	fmt.Println("token: ", token)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-type", "application/octet-stream")
	req.Header.Set("X-Goog-Upload-Content-Type", "video/mp4")
	req.Header.Set("X-Goog-Upload-Protocol", "raw")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	uploadToken := string(body)
	newItem := struct {
		AlbumId       *string        `json:"albumId,omitempty"`
		NewMediaItems []NewMediaItem `json:"newMediaItems"`
		AlbumPosition *AlbumPosition `json:"albumPosition"`
	}{
		NewMediaItems: []NewMediaItem{
			{
				Description: "",
				SimpleMediaItem: SimpleMediaItem{
					FileName:    fileName,
					UploadToken: uploadToken,
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(newItem)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonBytes))

	req, err = http.NewRequest("POST", "https://photoslibrary.googleapis.com/v1/mediaItems:batchCreate", bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d, response: %s", resp.StatusCode, resp.Status)
	}

	return nil

}
