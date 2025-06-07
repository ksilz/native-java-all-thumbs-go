package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

// Read image files from the current directory
func readImageFiles() ([]os.DirEntry, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	imageFiles := []os.DirEntry{}
	for _, entry := range entries {
		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				imageFiles = append(imageFiles, entry)
			}
		}
	}

	return imageFiles, nil
}

// Convert image to a one-page PDF
func imageToPDF(entry os.DirEntry) error {
	imgPath := entry.Name()
	pdfPath := filepath.Join("pdf", strings.TrimSuffix(imgPath, filepath.Ext(imgPath))+".pdf")

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set PDF metadata
	pdf.SetTitle("Converted image: "+imgPath, true)
	pdf.SetAuthor("Image to PDF Converter", true)

	// Add a new page
	pdf.AddPage()

	// Register the image
	imgOptions := gofpdf.ImageOptions{
		ImageType:            strings.TrimPrefix(filepath.Ext(imgPath), "."),
		ReadDpi:             true,
		AllowNegativePosition: false,
	}

	// Read the image file
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return fmt.Errorf("error opening image file: %w", err)
	}
	defer imgFile.Close()

	img, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return fmt.Errorf("error decoding image: %w", err)
	}

	// Calculate the scaling and positioning for the image in the PDF
	ratio := float64(img.Width) / float64(img.Height)
	margin := 10.0
	pageWidth, pageHeight := pdf.GetPageSize()
	width := pageWidth - 2*margin
	height := width / ratio
	x := margin
	y := (pageHeight - height) / 2

	// Place the image on the PDF
	pdf.ImageOptions(imgPath, x, y, width, height, false, imgOptions, 0, "")

	// Save the PDF
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		return fmt.Errorf("error saving PDF: %w", err)
	}

	return nil
}

func waitForEnter(message string) {
	fmt.Printf("\nPress ENTER %s...", message)
	var input string
	fmt.Scanln(&input)
}

const loopCount = 5

func main() {
	pid := os.Getpid()
	fmt.Println()
	fmt.Println("This program will convert all JPG and PNG pictures in the current directory into PDF.")
	fmt.Println()
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Println()

	fmt.Printf("Running with process ID: %d\n", pid)

	// Wait for a key press
	waitForEnter("to START")
	
	start := time.Now()

	// Read all image files in the current directory
	imageFiles, err := readImageFiles()
	if err != nil {
		fmt.Printf("Error reading image files: %v\n", err)
		return
	}

	if len(imageFiles) == 0 {
		fmt.Println("No image files found")
		return
	}

	for pass := 1; pass <= loopCount; pass++ {
		fmt.Printf("\nPass %d/%d\n", pass, loopCount)

		// Convert each image file to a one-page PDF
		for i, file := range imageFiles {
			fmt.Printf("\r  File %d", i+1)
			if err := imageToPDF(file); err != nil {
				fmt.Printf("\nError converting %s to PDF: %v\n", file.Name(), err)
			}
		}
	}

	stop := math.Round(float64(time.Since(start).Milliseconds()) / 10) / 100

	fmt.Printf("\nDone creating PDFs in %.1f seconds", stop)

	fmt.Println()
	waitForEnter("for garbage collection")
	
	fmt.Println("\nNow sleeping for 10 seconds, hoping for garbage collection.")
	
	runtime.GC()           // Force garbage collection
	runtime.Gosched()
	time.Sleep(10 * time.Second)
	
	fmt.Println("\nWoke up from sleep.")
	waitForEnter("to STOP")
	fmt.Println()
}
