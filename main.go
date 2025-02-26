package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go-cli-file-reader <directory-path> <output-md-file>")
		return
	}

	dirPath := os.Args[1]
	outputFile := os.Args[2]

	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("Directory '%s' does not exist.\n", dirPath)
		return
	}

	// Open or create the output .md file
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Printf("Failed to open or create the output file: %v\n", err)
		return
	}
	defer file.Close()

	// Walk through the directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Read the file content
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %v", path, err)
		}

		// Write the file path as a header
		header := fmt.Sprintf("\n# %s\n", path)
		if _, err := file.WriteString(header); err != nil {
			return fmt.Errorf("failed to write header for file '%s': %v", path, err)
		}

		// Write the opening backticks and "js"
		if _, err := file.WriteString("```js\n"); err != nil {
			return fmt.Errorf("failed to write opening backticks for file '%s': %v", path, err)
		}

		// Write the file content
		if _, err := file.Write(content); err != nil {
			return fmt.Errorf("failed to write content for file '%s': %v", path, err)
		}

		// Write the closing backticks
		if _, err := file.WriteString("\n```\n"); err != nil {
			return fmt.Errorf("failed to write closing backticks for file '%s': %v", path, err)
		}

		fmt.Printf("Processed file: %s\n", path)
		return nil
	})
	if err != nil {
		fmt.Printf("Error processing files: %v\n", err)
		return
	}

	fmt.Println("All files processed successfully and appended to the markdown file.")
}
