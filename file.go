package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Gitignore struct {
	Patterns []string
}

// Initialize ignore patterns from .gitignore
func NewGitignore(path string) *Gitignore {
	content, err := os.ReadFile(path)
	if err != nil {
		// If error reading .gitignore, assume no ignore patterns
		return &Gitignore{Patterns: []string{}}
	}
	lines := strings.Split(string(content), "\n")
	var ignorePatterns []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ignorePatterns = append(ignorePatterns, line)
	}
	return &Gitignore{Patterns: ignorePatterns}
}

// Check if a path matches any ignore pattern
func (g *Gitignore) IsIgnored(path string) bool {
	for _, pattern := range g.Patterns {
		// Handle directory patterns (ending with /)
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(path, dirPattern) {
				return true
			}
		}
		// Handle file patterns
		match, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && match {
			return true
		}
	}
	return false
}

func ListDirectory(path string, depth int) string {
	if depth <= 0 {
		return ""
	}

	path = Path(path)
	basePath := path

	// Look for .gitignore in the current directory
	gitignorePath := filepath.Join(path, ".gitignore")
	ignorePatterns := NewGitignore(gitignorePath)

	files, err := os.ReadDir(path)
	if err != nil {
		return ""
	}

	var fileNames []string
	for _, file := range files {
		relativePath := file.Name()
		// Skip ignored files
		if ignorePatterns.IsIgnored(relativePath) {
			continue
		}

		// Exclude hidden files and directories (starting with '.')
		if strings.HasPrefix(relativePath, ".") {
			continue
		}

		fileNames = append(fileNames, relativePath)

		if file.IsDir() && depth > 1 {
			subPath := filepath.Join(path, file.Name())
			subFiles := ListDirectory(subPath, depth-1)
			if subFiles != "" && subFiles != "Empty directory" {
				subFileList := strings.Split(subFiles, "\n")
				for _, subFile := range subFileList {
					// Get the relative path from the base path
					relPath, err := filepath.Rel(basePath, filepath.Join(subPath, subFile))
					if err == nil {
						fileNames = append(fileNames, relPath)
					}
				}
			}
		}
	}

	if len(fileNames) == 0 {
		return "Empty directory"
	}

	return strings.Join(fileNames, "\n")
}

func ReadFile(path string, offset int, length int) string {
	if length > 10000 {
		return "Cannot read more than 10000 lines"
	}

	if length == 0 {
		length = 10000
	}

	path = Path(path)

	// Reject if path points to a Go file
	if strings.HasSuffix(path, ".go") {
		return "Cannot read Go files directly. Use code functions instead."
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	text := string(content)
	lines := strings.Split(text, "\n")

	if offset > len(lines) {
		return fmt.Sprintf("File has %d lines, cannot read line %d", len(lines), offset)
	}

	var res string
	if offset+length > len(lines) {
		res = strings.Join(lines[offset:], "\n")
	} else {
		res = strings.Join(lines[offset:offset+length], "\n")
	}
	if res == "" {
		return "Empty file"
	}
	return res
}

func WriteFile(path string, content string) string {
	path = Path(path)

	// Reject if path points to a Go file
	if strings.HasSuffix(path, ".go") {
		return "Cannot write to Go files directly. Use code functions instead."
	}

	// Check if content contains partial patch markers
	hasPartialPatch := strings.Contains(content, "// rest of the code...") ||
		strings.Contains(content, "// ... existing code ...") ||
		strings.Contains(content, "// ...")

	var finalContent string
	if hasPartialPatch {
		// Read existing file content
		existingContent, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Sprintf("Error reading existing file: %v", err)
		}

		// If file doesn't exist, treat as new file
		if os.IsNotExist(err) {
			finalContent = content
		} else {
			// Split content into lines for processing
			existingLines := strings.Split(string(existingContent), "\n")
			newLines := strings.Split(content, "\n")
			var mergedLines []string

			// Process each line of new content
			for i := 0; i < len(newLines); i++ {
				line := newLines[i]
				if strings.Contains(line, "// rest of the code...") ||
					strings.Contains(line, "// ... existing code ...") ||
					strings.Contains(line, "// ...") {
					// When we hit a marker, include the rest of the existing content
					mergedLines = append(mergedLines, existingLines...)
					break
				}
				mergedLines = append(mergedLines, line)
			}

			// Join the lines back together
			finalContent = strings.Join(mergedLines, "\n")
		}
	} else {
		finalContent = content
	}

	err := os.WriteFile(path, []byte(finalContent), 0644)
	if err != nil {
		return fmt.Sprintf("Error writing to file: %v", err)
	}

	lint := Lint(path)

	if strings.Contains(lint, "no Go files") {
		lint = ""
	}

	return fmt.Sprintf("Path: %s\n\nNew content:\n%s\n\n---\n\nLinter results:\n%s", path, finalContent, lint)
}

func MkDir(path string) string {
	path = Path(path)

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}
	return fmt.Sprintf("Directory created: %s", path)
}

func FetchWikiDocs() string {
	wikiPath := "wiki"
	files, err := os.ReadDir(wikiPath)
	if err != nil {
		return fmt.Sprintf("Error reading wiki directory: %v", err)
	}

	var docs []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			content, err := os.ReadFile(filepath.Join(wikiPath, file.Name()))
			if err != nil {
				docs = append(docs, fmt.Sprintf("Error reading %s: %v", file.Name(), err))
				continue
			}
			docs = append(docs, fmt.Sprintf("File: %s\n\nContent:\n%s", file.Name(), content))
		}
	}
	if len(docs) == 0 {
		return "No Markdown files found in wiki."
	}
	return strings.Join(docs, "\n\n---\n\n")
}

func SearchText(query string) string {
	return searchTextRecursive(workingDirectory, query)
}

func searchTextRecursive(dir string, query string) string {
	dir = Path(dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err)
	}

	var results []string
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		if file.IsDir() {
			// Recursively search subdirectories
			subResults := searchTextRecursive(filePath, query)
			if subResults != "No results found" {
				results = append(results, subResults)
			}
			continue
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		text := string(content)
		lines := strings.Split(text, "\n")

		for i, line := range lines {
			if strings.Contains(line, query) {
				// Get relative path from working directory
				relPath, _ := filepath.Rel(workingDirectory, filePath)
				results = append(results, fmt.Sprintf("%s:%d: %s", relPath, i+1, line))
			}
		}
	}

	if len(results) == 0 {
		return "No results found"
	}
	return strings.Join(results, "\n")
}

func Path(path string) string {
	if path == "." || path == "" {
		path = workingDirectory
	}
	// if relative path, convert to absolute path
	if !filepath.IsAbs(path) {
		path = filepath.Join(workingDirectory, path)
	}
	return path
}
