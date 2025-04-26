package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ListDirectory(path string, depth int) []string {
	if depth < 0 {
		return []string{}
	}

	if path == "." || path == "" {
		path = workingDirectory
	}

	// if relative path, convert to absolute path
	if !filepath.IsAbs(path) {
		path = filepath.Join(workingDirectory, path)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return []string{}
	}

	var fileNames []string
	for _, file := range files {
		relativePath := file.Name()
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

	os.WriteFile(path, []byte(strings.Join(lines[:offset], "\n")), 0644)
	os.WriteFile(path, []byte(strings.Join(lines[offset:], "\n")), 0644)

	return fmt.Sprintf("Wrote %d lines to %s", len(lines), path)
}
