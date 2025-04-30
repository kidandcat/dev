package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client
var logger *log.Logger

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

	logFilePath := filepath.Join(workingDirectory, "LOG.txt")
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s", err)
		os.Exit(1)
	}
	logger = log.New(logFile, "", 0)

	client = openai.NewClient("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")

	// check that a file INPUT.md exists in the working directory
	if _, err := os.Stat(filepath.Join(workingDirectory, "INPUT.md")); os.IsNotExist(err) {
		fmt.Printf("Input file INPUT.md does not exist in the working directory %s", workingDirectory)
		os.Exit(1)
	}

	for {
		handleChatCompletion(MODEL_GPT41, openai.ChatCompletionMessage{
			Role:    "user",
			Content: "continue",
		})
	}
}
