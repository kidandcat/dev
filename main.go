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
		workingDirectory = "."
	} else {
		workingDirectory = os.Args[1]
	}
	if _, err := os.Stat(workingDirectory); os.IsNotExist(err) {
		fmt.Printf("Working directory %s does not exist", workingDirectory)
		os.Exit(1)
	}

	if !filepath.IsAbs(workingDirectory) {
		var err error
		workingDirectory, err = filepath.Abs(workingDirectory)
		if err != nil {
			fmt.Printf("Error converting working directory to absolute path: %s", err)
			os.Exit(1)
		}
	}

	client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	if _, err := os.Stat(filepath.Join(workingDirectory, "INPUT.md")); os.IsNotExist(err) {
		fmt.Printf("Input file INPUT.md does not exist in the working directory %s", workingDirectory)
		os.Exit(1)
	}

	handleChatCompletion(MODEL_GPT41, openai.ChatCompletionMessage{
		Role: "user",
		Content: `
		Open a file called INPUT.md and read the content.
		Add new tasks based on the content of the INPUT.md file to the markdown checklist in the TASKS.md file.
		
		If the file TASKS.md does not exist, create it.
		`,
	})
	for {
		response := handleChatCompletion(MODEL_GPT41, openai.ChatCompletionMessage{
			Role: "user",
			Content: `
			Read the TASKS.md file and do the next task.
			After each task, update the TASKS.md file to reflect the changes.
			`,
		})
		if YesNoQuestion(fmt.Sprintf("Has all the tasks been completed? %s", response)) {

			inputPath := filepath.Join(workingDirectory, "INPUT.md")
			err := os.WriteFile(inputPath, []byte{}, 0644)
			if err != nil {
				fmt.Printf("Error erasing INPUT.md: %s", err)
			}
			break
		}
	}
}
