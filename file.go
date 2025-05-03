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

func ListDirectory(path string, depth int) map[string]any {
	if depth < 0 {
		return map[string]any{
			"error": "Depth cannot be negative",
		}
	}

	path = Path(path)

	// Look for .gitignore in the current directory
	gitignorePath := filepath.Join(path, ".gitignore")
	ignorePatterns := NewGitignore(gitignorePath)

	files, err := os.ReadDir(path)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error reading directory: %v", err),
		}
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
			if subFiles["error"] != "Empty directory" {
				if subFiles["files"] != nil {
					subFileList := subFiles["files"].([]string)
					for _, subFile := range subFileList {
						fileNames = append(fileNames, filepath.Join(relativePath, subFile))
					}
				}
			}
		}
	}

	if len(fileNames) == 0 {
		return map[string]any{
			"error": "Empty directory",
		}
	}

	return map[string]any{
		"files": fileNames,
	}
}

func ReadFile(path string, offset int, length int) map[string]any {
	if length > 1000 {
		return map[string]any{
			"error": "Cannot read more than 1000 lines",
		}
	}

	if length == 0 {
		length = 1000
	}

	path = Path(path)

	// Reject if path points to a Go file
	if strings.HasSuffix(path, ".go") {
		return map[string]any{
			"error": "Cannot read Go files directly. Use code functions instead.",
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error reading file: %v", err),
		}
	}
	text := string(content)
	lines := strings.Split(text, "\n")

	if offset > len(lines) {
		return map[string]any{
			"error": fmt.Sprintf("File has %d lines, cannot read line %d", len(lines), offset),
		}
	}

	var res string
	if offset+length > len(lines) {
		res = strings.Join(lines[offset:], "\n")
	} else {
		res = strings.Join(lines[offset:offset+length], "\n")
	}
	if res == "" {
		return map[string]any{
			"error": "Empty file",
		}
	}
	return map[string]any{
		"content": res,
	}
}

func WriteFile(path string, content string) map[string]any {
	path = Path(path)

	// Check if content contains partial patch markers
	hasPartialPatch := strings.Contains(content, "// rest of the code...") ||
		strings.Contains(content, "// ... existing code ...") ||
		strings.Contains(content, "// ...")

	var finalContent string
	if hasPartialPatch {
		// Read existing file content
		existingContent, err := os.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return map[string]any{
				"error": fmt.Sprintf("Error reading existing file: %v", err),
			}
		}

		// If file doesn't exist, treat as new file
		if os.IsNotExist(err) {
			finalContent = content
		} else {
			// Reject if path points to a Go file
			if strings.HasSuffix(path, ".go") {
				return map[string]any{
					"error": "Cannot write to existing Go files directly. Use code functions instead.",
				}
			}

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
		return map[string]any{
			"error": fmt.Sprintf("Error writing to file: %v", err),
		}
	}

	return map[string]any{
		"path":    path,
		"content": finalContent,
		"lint":    Lint(path),
	}
}

func MkDir(path string) map[string]any {
	path = Path(path)

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error creating directory: %v", err),
		}
	}
	return map[string]any{
		"path": path,
	}
}

func FetchWikiDocs() map[string]any {
	wikiPath := "wiki"
	files, err := os.ReadDir(wikiPath)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error reading wiki directory: %v", err),
		}
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
		return map[string]any{
			"error": "No Markdown files found in wiki.",
		}
	}
	return map[string]any{
		"results": strings.Join(docs, "\n\n---\n\n"),
	}
}

func SearchText(query string) map[string]any {
	return searchTextRecursive(workingDirectory, query)
}

func searchTextRecursive(dir string, query string) map[string]any {
	dir = Path(dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error reading directory: %v", err),
		}
	}

	var results []string
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		if file.IsDir() {
			// Recursively search subdirectories
			subResults := searchTextRecursive(filePath, query)
			if subResults["error"] != "No results found" {
				if subResults["results"] != nil {
					subResultsList := subResults["results"].([]string)
					results = append(results, subResultsList...)
				}
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
		return map[string]any{
			"error": "No results found",
		}
	}
	return map[string]any{
		"results": results,
	}
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
