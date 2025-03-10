package main

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

func main() {

	token := "Bearer ya29.a0AeXRPp74Mw2rez1cCFtxqFBkq1L3eoIcTZln1fvGhiO30ngOMHH1PEVp39LhUqKILfj_yDmsafkI0dqYgf8vcCOQbye10eiYhAqvjWkf19v6WwQTMACuQtZ7SsRLp6Y1G9YdVashqla4Xa_KPC8uO7Afc-xea-sF-30_qByzaCgYKAb0SARESFQHGX2MiwnXlqLVapdSn-e2Jwms8KQ0175"

	filename := "6a3cddbf-11a0-43d0-9f35-be1ce2490807"

	filePath := fmt.Sprintf("/tmp/out/%s.mp4", filename)

	fmt.Println("Uploading video...")

	err := uploadVideo(token, filePath, "PXL_20250227_233455922.mp4")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Video uploaded successfully!")
	}
}

func uploadVideo(token string, videoPath string, fileName string) error {
	file, err := os.Open(videoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("POST", "https://photoslibrary.googleapis.com/v1/uploads", file)
	if err != nil {
		return err
	}

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
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Errorf("status code: %d", resp.StatusCode)
			return err
		}
		return fmt.Errorf("status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil

}
