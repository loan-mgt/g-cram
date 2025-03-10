package main

import (
	"fmt"
	"os/exec"
)

func main() {

	filename := "6a3cddbf-11a0-43d0-9f35-be1ce2490807"
	inputPath := fmt.Sprintf("/tmp/in/%s.mp4", filename)
	outputPath := fmt.Sprintf("/tmp/out/%s.mp4", filename)

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx265", "-crf", "28", outputPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", output)
		return
	}

	fmt.Println("Conversion successful!")
}
