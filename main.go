package main

import (
	"fmt"
	"image"
	_ "github.com/chai2010/webp" // Import the package to support 16-bit depth PNG
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	return err
}

func waitForEnter() {
	fmt.Println("Press ENTER to continue...")
	var input string
	fmt.Scanln(&input)
}

func main() {

	pid := os.Getpid()
	fmt.Printf("Running with process ID: %d\n", pid)

	// Wait for a key press
	waitForEnter()
	
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

	// Create the "pdf" sub-directory if it doesn't exist
	if _, err := os.Stat("pdf"); os.IsNotExist(err) {
		os.Mkdir("pdf", 0755)
	}

	// Convert each image file to a one-page PDF
	for i, file := range imageFiles {
		fmt.Printf("Processing file %d: %s\n", i+1, file.Name())
		err := imageToPDF(file)
		if err != nil {
			if strings.Contains(err.Error(), "16-bit depth not supported") {
				fmt.Printf("Skipped %s: %v\n", file.Name(), err)
			} else {
				fmt.Printf("Error converting %s to PDF: %v\n", file.Name(), err)
			}
		} else {
			fmt.Printf("Successfully converted %s to PDF\n", file.Name())
		}
	}
}
