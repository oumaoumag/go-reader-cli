package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/v5/plumbing"
)

// ShouldSkipDir determines if a directory should be skipped based on .gitignore patterns.
func shouldSkipDir(relPath string, patterns []string) bool {
	// Esxplicitly skip node_modules
	if  strings.Contains(relPath, "lib") || strings.Contains(relPath, "node_modules")|| strings.Contains(relPath, "vender") || strings.Contains(relPath, "vender") || strings.Contains(relPath, "tmp") || strings.Contains(relPath, "log") {
		return true
	}
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "/") || strings.HasPrefix(pattern, "/")  {
			dirPattern := strings.TrimSuffix(pattern, "/")
			dirPattern = strings.TrimPrefix(dirPattern, "/")
			if match, _ := filepath.Match(dirPattern, relPath); match {
				return true
			}
		} else {
			if match, _ := filepath.Match(pattern, relPath); match {
				return true
			}
		}
	}
	return false
}

// shouldSkipFile determines if a file should be skipped based on .gitignore patterns.
func shouldSkipFile(relPath string, name string, patterns []string) bool {
	// Skip all files in the bin/directory
	if strings.HasPrefix(relPath, "bin/") {
		return true
	}

	// Skip specific SQLite database files
	if name == "development.sqlite3" {
		return true
	}
	for _, pattern := range patterns {
		if !strings.HasSuffix(pattern, "/") || strings.HasPrefix(pattern, "/"){
			if strings.Contains(pattern, "/") {
				if match, _ := filepath.Match(pattern, relPath); match {
					return true
				}
			} else {
				if match, _ := filepath.Match(pattern, name); match {
					return true
				}
			}
		}
	}
	return false
}

func checkIfGitOrFilePath(input string) (repoURL string, branch string, isRepo bool) {
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")  || strings.HasPrefix(input, "git@") {
		parts := strings.Split(input, "@")
		if strings.HasPrefix(input, "git@") {
			// SSH format: git@github.com:user/repo.git@branch
			if len(parts) > 2 {
				repoURL = strings.Join(parts[:2], "@")
				branch = parts[2]
			} else {
				repoURL = input
			}
		} else {
			// HTTPS format: https://github.com/user/repo.git
			if len(parts) > 1 {
				repoURL = parts[0]
				branch = parts[1]
			}  else {
				repoURL = input
			}
		}
		return repoURL, branch, true
	}
	return "", "", false
}

func cloneRepo(repoURL, branch string) (string, error) {
	tempDir, err := os.MkdirTemp("",  "go-reader-cli-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}

	cloneOptions := &git.cloneOptions{
		URL: repoURL,
	}

	if branch != "" {
		cloneOptions.ReferenceName  = plumbing.ReferenceName("refs/heads/" + branch)
		cloneOptions.SingleBranch = true
	}

	if branch != "" {
		cloneOptions.ReferenceName = plumbing.ReferenceName("refs/heads/" + branch)
		cloneOptions.SingleBranch = true
	}

	_, err = git.PlainClone(tempDir, false, cloneOptions)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone repository: %v", err)
	}
	return tempDir, nil                                                                                                                                                                                                                                                                                                                                                                                                       
}

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

	// Extensions to skip (images and config files)
	skipExtensions := map[string]struct{} {
		".jpg": {},
		".jpeg": {},
		".png": {},
		".gif": {},
		".bmp": {},
		"tiff": {},
		".svg": {},
		".mp3": {},
		".mp4": {},
		// ".config": {},
		// ".config.ts": {},
		".config.mjs": {},
		".ini": {},
		".yaml": {},
		".yml": {},
		".toml": {},
		".ico": {},
		".sqlite3": {},
		".db": {},
		// ".json": {},
		}

		// Read .gitigone if it exists
		patterns := []string{}
		gitignorePath := filepath.Join(dirPath, ".gitignore")
		if _, err := os.Stat(gitignorePath); err == nil {
			content, err := os.ReadFile(gitignorePath)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && !strings.HasPrefix(line, "#") {
						patterns = append(patterns, line)
					}
				}
			}
		}

	// Walk through the directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			fmt.Printf("Failed to get the relative path for %s: %v\n", path, err)
			relPath = path
		}

		// Skip directories matching .gitignore patterns or starting with a dot
		if info.IsDir() {
			if (path != dirPath && strings.HasPrefix(info.Name(), ".")) || shouldSkipDir(relPath, patterns) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip files matching .gitignore patterns, starting with a .dot, in bin/, or with skipped extensions
		if shouldSkipFile(relPath, info.Name(), patterns) || strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := skipExtensions[ext]; ok {
			return nil
		}

		// Read the file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file '%s': %v", path, err)
		}

		// Write the file path as a header using relative Path
		header := fmt.Sprintf("\n# %s\n", relPath)
		if _, err := file.WriteString(header); err != nil {
			return fmt.Errorf("failed to write header for file '%s': %v", path, err)
		}

		// Write the opening backticks with appropirate language identifier
		lang := getLanguageIdentifier(ext)

		if _, err := file.WriteString(fmt.Sprintf("```%s\n", lang)); err != nil {
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

// getLanguageIdentifier returns the appropriate Markdown code block language base on the file extentions
func getLanguageIdentifier(ext string) string {
	switch strings.ToLower(ext) {
	case ".rb":
		return "ruby"
	case ".js":
		return "javascript"
	case ".scss":
		return "scss"
	case ".html", "erb":
		return "html"
	case ".md":
		return "markdown"
	case ".yml", ".yaml":
		return "yaml"
	case ".rs":
		return "rust"
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".sol":
		return "solidity"
	default:
		return ""
		
	}
}