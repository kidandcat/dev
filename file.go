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
		match, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && match {
			return true
		}
	}
	return false
}

func ListDirectory(path string, depth int) []string {
	if depth < 0 {
		return []string{}
	}

	// Initialize ignore patterns if not already done
	ignorePatterns := NewGitignore(path)

	path = Path(path)

	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}
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

		if file.IsDir() && depth > 0 {
			subPath := filepath.Join(path, file.Name())
			subFiles := ListDirectory(subPath, depth-1)
			for _, subFile := range subFiles {
				fileNames = append(fileNames, filepath.Join(relativePath, subFile))
			}
		}
	}

	return fileNames
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
		return ""
	}
	text := string(content)
	lines := strings.Split(text, "\n")

	if offset > len(lines) {
		return fmt.Sprintf("File has %d lines, cannot read line %d", len(lines), offset)
	}

	if offset+length > len(lines) {
		return strings.Join(lines[offset:], "\n")
	}

	return strings.Join(lines[offset:offset+length], "\n")
}

func WriteFile(path string, content string, offset int) string {
	content = strings.ReplaceAll(content, "\r", "")
	lines := strings.Split(content, "\n")

	if offset > len(lines) {
		offset = len(lines)
	}

	path = Path(path)

	os.WriteFile(path, []byte(strings.Join(lines[:offset], "\n")), 0644)
	os.WriteFile(path, []byte(strings.Join(lines[offset:], "\n")), 0644)

	res := fmt.Sprintf("Wrote %d lines to %s", len(lines), path)
	lint := Lint(path)

	return res + "\n" + lint
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
