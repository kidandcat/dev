package main

import (
	"fmt"
	"os"
	"path/filepath"

	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <working directory>")
		os.Exit(1)
	}
	workingDirectory = os.Args[1]
	if _, err := os.Stat(workingDirectory); os.IsNotExist(err) {
		fmt.Printf("Working directory %s does not exist", workingDirectory)
		os.Exit(1)
	}
	// if working directory is not absolute, convert to absolute path
	if !filepath.IsAbs(workingDirectory) {
		var err error
		workingDirectory, err = filepath.Abs(workingDirectory)
		if err != nil {
			fmt.Printf("Error converting working directory to absolute path: %s", err)
			os.Exit(1)
		}
	}

	client = openai.NewClient("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")

	model := NewModel()
	model.Start()
}
