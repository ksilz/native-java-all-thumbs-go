package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

// Read image files from the current directory
func readImageFiles() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}

	imageFiles := []os.FileInfo{}
	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				imageFiles = append(imageFiles, file)
			}
		}
	}

	return imageFiles, nil
}

// Convert image to a one-page PDF
func imageToPDF(file os.FileInfo) error {
	imgPath := file.Name()
	pdfPath := "pdf/" + strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + ".pdf"

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set PDF metadata
	pdf.SetTitle("Converted image: " + imgPath, true)
	pdf.SetAuthor("Image to PDF Converter", true)

	// Add a new page
	pdf.AddPage()

	// Register the image
	imgOptions := gofpdf.ImageOptions{ImageType: strings.TrimPrefix(filepath.Ext(imgPath), ".")}

	// Read the image file using an io.Reader
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	img, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return err
	}

	// Calculate the scaling and positioning for the image in the PDF
	ratio := float64(img.Width) / float64(img.Height)
	margin := 10.0
	width, height := pdf.GetPageSize()
	width -= 2 * margin
	height = width / ratio
	x := margin
	y := (height - height) / 2

	// Place the image on the PDF
	pdf.ImageOptions(imgPath, x, y, width, height, false, imgOptions, 0, "")

	// Save the PDF
	err = pdf.OutputFileAndClose(pdfPath)

	// Reduce memory usage
	pdf = nil              // Set pdf to nil to release the memory

	return err
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
	
	var start = time.Now()

	// Read all image files in the current directory
	imageFiles, err := readImageFiles()
	if err != nil {
		fmt.Println("Error reading image files:", err)
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
			err := imageToPDF(file)
			if err != nil {
					fmt.Printf("Error converting %s to PDF: %v\n", file.Name(), err)
			}
		}
	}

  var stop = math.Round(float64(time.Since(start).Milliseconds()) / 10) / 100

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
