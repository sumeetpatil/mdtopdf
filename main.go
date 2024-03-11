package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jung-kurt/gofpdf"
	"github.com/russross/blackfriday/v2"
)

// stripHTMLTags removes HTML tags from the given input.
func stripHTMLTags(input string) string {
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(input, "")
}

func main() {
	// Define command-line flags
	inputDir := flag.String("input", ".", "Input directory containing .md files")
	outputPDF := flag.String("output", "merged_output.pdf", "Output PDF file")
	flag.Parse()

	// Slice to store all .md files
	var files []string

	// Walk through the directory and its subdirectories
	err := filepath.Walk(*inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		// Check if the file has a .md extension
		if filepath.Ext(path) == ".md" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directory:", err)
		os.Exit(1)
	}

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set the font for Markdown content
	pdf.SetFont("Arial", "", 12)

	// Iterate through each .md file
	for _, file := range files {
		// Print the file being processed
		fmt.Printf("Processing file: %s\n", file)

		// Read the content of the .md file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file, err)
			continue
		}

		// Convert Markdown to HTML
		htmlContent := string(blackfriday.Run(content))

		// Strip HTML tags
		plainContent := stripHTMLTags(htmlContent)

		// Add the plain content to the PDF document
		pdf.AddPage()
		pdf.MultiCell(0, 10, plainContent, "", "L", false)
	}

	// Output the merged PDF to the specified file
	err = pdf.OutputFileAndClose(*outputPDF)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Merged .md files into %s\n", *outputPDF)
}

