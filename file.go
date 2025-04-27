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
	if depth < 0 {
		return ""
	}

	path = Path(path)

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
					fileNames = append(fileNames, filepath.Join(relativePath, subFile))
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
	if length > 1000 {
		return fmt.Sprintf("Cannot read more than 1000 lines")
	}

	if length == 0 {
		length = 1000
	}

	path = Path(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	text := string(content)
	lines := strings.Split(text, "\n")

	if offset > len(lines) {
		return fmt.Sprintf("File has %d lines, cannot read line %d", len(lines), offset)
	}

	if offset+length > len(lines) {
		return strings.Join(lines[offset:], "\n")
	}

	res := strings.Join(lines[offset:offset+length], "\n")
	if res == "" {
		return "Empty file"
	}
	return res
}

func WriteFile(path string, content string) string {
	path = Path(path)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return fmt.Sprintf("Error writing to file: %v", err)
	}

	lint := Lint(path)

	if strings.Contains(lint, "no Go files") {
		lint = ""
	}

	return fmt.Sprintf("Path: %s\n\nNew content:\n%s\n\n---\n\nLinter results:\n%s", path, content, lint)
}

func MkDir(path string) string {
	path = Path(path)

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}
	return fmt.Sprintf("Directory created: %s", path)
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
