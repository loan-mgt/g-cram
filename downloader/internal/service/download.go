package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func DownloadVideo(msg Msg) (int64, error) {

	if !isGoogleStorageURL(msg.BaseUrl) {
		return 0, fmt.Errorf("Invalid BaseURL")
	}

	baseURL := msg.BaseUrl + "=dv"

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", msg.Token)

	fileSize, err := saveVideoToFile(msg, req)
	if err != nil {
		return 0, err
	}

	return fileSize, nil
}

func saveVideoToFile(msg Msg, req *http.Request) (int64, error) {
	videoPath := fmt.Sprintf("/tmp/in/%s_%d/%s.mp4", msg.UserId, msg.Timestamp, msg.MediaId)

	// create the directory if it doesn't exist
	err := os.MkdirAll(fmt.Sprintf("/tmp/in/%s_%d", msg.UserId, msg.Timestamp), 0755)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	fmt.Println(videoPath)
	f, err := os.Create(videoPath)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer f.Close()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return resp.ContentLength, nil
}

func isGoogleStorageURL(url string) bool {
	return regexp.MustCompile(`^https:\/\/[a-zA-Z0-9-]*\.googleusercontent\.com\/.*$`).MatchString(url)
}
