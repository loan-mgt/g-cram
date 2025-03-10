package service

import (
	"fmt"
	"os/exec"
)

func CreateVideoWithMetadata(filename, creationDate string) error {
	inputPath := fmt.Sprintf("/tmp/in/%s.mp4", filename)
	outputPath := fmt.Sprintf("/tmp/out/%s.mp4", filename)

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx265", "-crf", "28", "-metadata",
		fmt.Sprintf("creation_time=%s", creationDate), "-c:v", "copy", "-c:a", "copy", "-f", "mp4", outputPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", output)
		return err
	}

	return nil
}
