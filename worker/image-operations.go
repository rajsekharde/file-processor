package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

/*
Opens and decodes input file, resizes to given dimensions and saves as a PNG image.
Returns (0, output filename) or (1, Error message) if operation fails.
*/
func resizeImage(taskID string, inputFileName string, targetWidth int, targetHeight int) (int, string) {
	inputPath := "../file-storage/uploads/" + inputFileName
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return 1, err.Error()
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return 1, err.Error()
	}
	
	outputFileName := "output-" + taskID + ".png"
	outputPath := "../file-storage/completed/" + outputFileName

	// Create a new blank RGBA destination image
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))

	// Scale using BiLinear interpolation (CatmullRom is also available for higher quality)
	draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return 1, err.Error()
	}
	defer outputFile.Close()

	// encode the image into the output file as PNG
	if err := png.Encode(outputFile, dst); err != nil {
		return 1, err.Error()
	}

	return 0, outputFileName
}

/*
Opens and converts input file to given format, and stores the output file.
Returns (0, output filename) or (1, Error message) if operation fails.
*/
func convertFormat(taskID string, inputFileName string, outputFormat string) (int, string) {
	outputFileName := "output-" + taskID + "." + outputFormat
	outputPath := "../file-storage/completed/" + outputFileName
	inputPath := "../file-storage/uploads/" + inputFileName

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

	return 0, outputFileName
}