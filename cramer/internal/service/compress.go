package service

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func CreateVideoWithMetadata(filename string, creationDate int64) (int64, error) {
	inputPath := fmt.Sprintf("/tmp/in/%s.mp4", filename)
	outputPath := fmt.Sprintf("/tmp/out/%s.mp4", filename)

	if _, err := os.Stat(outputPath); err == nil {
		err = os.Remove(outputPath)
		if err != nil {
			return 0, err
		}
	}

	creationDateFormated := time.Unix(0, creationDate*int64(time.Millisecond)).Format("2006-01-02T15:04:05.000Z")

	// log creation date
	fmt.Println("Creation date:", creationDate, creationDateFormated)

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx265", "-crf", "28", "-metadata",
		fmt.Sprintf("creation_time=%s", creationDateFormated), "-c:a", "copy", "-f", "mp4", outputPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", output)
		return 0, err
	}

	inputStat, err := os.Stat(inputPath)
	if err != nil {
		return 0, err
	}
	outputStat, err := os.Stat(outputPath)
	if err != nil {
		return 0, err
	}

	inputSize := inputStat.Size()
	outputSize := outputStat.Size()
	percent := float64(outputSize) / float64(inputSize) * 100
	fmt.Printf("Ouput size: %d%% (%d bytes -> %d bytes)\n", int(percent), inputSize, outputSize)

	return outputSize, nil
}
