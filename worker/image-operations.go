package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func ConvertFormat(inputPath string, outputFormat string) int {
	outputPath := "../file-storage/completed/output." + outputFormat

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return 1
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return 1
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating destination file: %v\n", err)
		return 1
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
		return 1
	}
	if err != nil {
		fmt.Printf("Error encoding image: %v\n", err)
		return 1
	}

	fmt.Printf("Image successfully converted and saved\n")
	return 0
}