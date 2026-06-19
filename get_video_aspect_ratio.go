package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffprobe error: %w", err)
	}

	var output struct {
		Streams []struct {
			Width              int    `json:"width"`
			Height             int    `json:"height"`
			DisplayAspectRatio string `json:"display_aspect_ratio"`
		} `json:"streams"`
	}

	err = json.Unmarshal(stdout.Bytes(), &output)
	if err != nil {
		return "", fmt.Errorf("couldn't parse ffprobe output: %w", err)
	}

	for _, stream := range output.Streams {
		switch stream.DisplayAspectRatio {
		case "16:9":
			return "landscape", nil
		case "9:16":
			return "portrait", nil
		}

		if stream.Width > 0 && stream.Height > 0 {
			ratio := float64(stream.Width) / float64(stream.Height)
			landscape := 16.0 / 9.0 // ~1.778
			portrait := 9.0 / 16.0  // ~0.5625

			tolerance := 0.05

			if math.Abs(ratio-landscape) < tolerance {
				return "landscape", nil
			} else if math.Abs(ratio-portrait) < tolerance {
				return "portrait", nil
			}
		}
	}

	return "other", nil
}
