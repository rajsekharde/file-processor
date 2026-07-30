package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ConvertFormat(inputPath string, outputFormat string) (int, string) {
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")
	outputPath := "../file-storage/completed/output_" + dateString + "." + outputFormat

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return 1, "Error opening input file"
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return 1, "Error decoding image"
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating destination file: %v\n", err)
		return 1, "Error creating destination file"
	}
	defer outputFile.Close()

	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outputFile, img, nil)
	case ".png":
		err = png.Encode(outputFile, img)
	default:
		fmt.Printf("Unsupported output format: %s\n", ext)
		return 1, "Unsupported output format"
	}
	if err != nil {
		fmt.Printf("Error encoding image: %v\n", err)
		return 1, "Error encoding image"
	}

	fmt.Printf("Image successfully converted and saved\n")
	return 0, outputPath
}