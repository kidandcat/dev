package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/genai"
)

var client *genai.Client

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

	key := os.Getenv("GEMINI_API_KEY")
	if key == "" {
		fmt.Printf("GEMINI_API_KEY is not set")
		os.Exit(1)
	}

	var err error
	ctx := context.Background()
	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		fmt.Printf("Error creating client: %s", err)
		os.Exit(1)
	}

	if _, err := os.Stat(filepath.Join(workingDirectory, "INPUT.md")); os.IsNotExist(err) {
		fmt.Printf("Input file INPUT.md does not exist in the working directory %s", workingDirectory)
		os.Exit(1)
	}

	handleChatCompletion(MODEL_BIG, &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			genai.NewPartFromText(`
			Open a file called INPUT.md and read the content.
			Split the content of the INPUT.md file into tasks and add them to the markdown checklist in the TASKS.md file.

			Follow this template:
			- [ ] Task 1: Detailed description of the task.
				- [ ] Implementation
				- [ ] Tests
			- [ ] Task 2: Detailed description of the task.
				- [ ] Implementation
				- [ ] Tests

			If the file TASKS.md does not exist, create it.
			`),
		},
	})

	for {
		response := handleChatCompletion(MODEL_BIG, &genai.Content{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				genai.NewPartFromText(`
				Read the TASKS.md file and do the next task.
				After each task, update the TASKS.md file to reflect the changes.

				Do not use placeholders, todo, etc. The task must be completed, including tests.
				`),
			},
		})
		tasks, err := os.ReadFile(filepath.Join(workingDirectory, "TASKS.md"))
		if err != nil {
			fmt.Printf("Error reading TASKS.md: %s", err)
			os.Exit(1)
		}
		if YesNoQuestion(fmt.Sprintf(`Has all the tasks been completed?
		
		Tasks:
		%s
		
		Response:
		%s
		
		`, string(tasks), response)) {
			inputPath := filepath.Join(workingDirectory, "INPUT.md")
			err := os.WriteFile(inputPath, []byte{}, 0644)
			if err != nil {
				fmt.Printf("Error erasing INPUT.md: %s", err)
			}
			break
		}
	}
}
