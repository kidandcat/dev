package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
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

	key := os.Getenv("OPENROUTER_API_KEY")
	if key == "" {
		fmt.Printf("OPENROUTER_API_KEY is not set")
		os.Exit(1)
	}
	config := openai.DefaultConfig(key)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client = openai.NewClientWithConfig(config)

	if _, err := os.Stat(filepath.Join(workingDirectory, "INPUT.md")); os.IsNotExist(err) {
		fmt.Printf("Input file INPUT.md does not exist in the working directory %s", workingDirectory)
		os.Create(filepath.Join(workingDirectory, "INPUT.md"))
		os.Exit(1)
	}

	if _, err := os.Stat(filepath.Join(workingDirectory, "TASKS.md")); os.IsNotExist(err) {
		os.Create(filepath.Join(workingDirectory, "TASKS.md"))
	}

	// GenWiki()

	handleChatCompletion(openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		Content: `
			Open a file called INPUT.md and read the content.
			Process the content of the INPUT.md file into independent, small tasks and add them to the markdown checklist in the TASKS.md file.

			If the file TASKS.md does not exist, create it.
		`,
	})

	for {
		messages = nil
		tasks, err := os.ReadFile(filepath.Join(workingDirectory, "TASKS.md"))
		if err != nil {
			fmt.Printf("Error reading TASKS.md: %s", err)
			os.Exit(1)
		}
		response := handleChatCompletion(openai.ChatCompletionMessage{
			Role: openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(`
				Do the next task.
				After each task, update the TASKS.md file to reflect the changes.

				TASKS:
				%s
			`, string(tasks)),
		})
		if response == "no_response" {
			log.Printf("No response from assistant, finishing")
			break
		}
		if !YesNoQuestion(fmt.Sprintf(`Has all the tasks been completed?
		
		Tasks:
		%s
		
		Response:
		%s
		
		`, string(tasks), response)) {
			log.Printf("Tasks not completed, continuing")
			continue
		}
		if ArePendingTodos() {
			diff, err := exec.Command("git", "diff").Output()
			if err != nil {
				fmt.Printf("Error running git diff: %s", err)
				os.Exit(1)
			}
			handleChatCompletion(openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleUser,
				Content: fmt.Sprintf(`
					Create tasks in the TASKS.md file to implement the missing functionality based on the TODOs, placeholders, etc. in the following git diff:

					git diff:
					%s
					`, string(diff)),
			})
			log.Printf("There are pending todos, continuing")
			continue
		}
		// Erase the INPUT.md file
		if err := os.WriteFile(filepath.Join(workingDirectory, "INPUT.md"), []byte{}, 0644); err != nil {
			fmt.Printf("Error erasing INPUT.md: %s", err)
		}
		break
	}
}

func ArePendingTodos() bool {
	diff, err := exec.Command("git", "diff").Output()
	if err != nil {
		fmt.Printf("Error running git diff: %s", err)
		return false
	}

	return YesNoQuestion(fmt.Sprintf(`
	Check if there are any TODOs, placeholders, etc. in the following git diff:
	
	%s
	`, string(diff)))
}
