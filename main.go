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

	client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	if _, err := os.Stat(filepath.Join(workingDirectory, "INPUT.md")); os.IsNotExist(err) {
		fmt.Printf("Input file INPUT.md does not exist in the working directory %s", workingDirectory)
		os.Exit(1)
	}

	handleChatCompletion(MODEL_BIG, openai.ChatCompletionMessage{
		Role:    "user",
		Content: "Open a file called INPUT.md and read the content. Process the content of the INPUT.md file into independent, small tasks and add them to the markdown checklist in the TASKS.md file.",
	})

	for {
		response := handleChatCompletion(MODEL_SMALL, openai.ChatCompletionMessage{
			Role: "user",
			Content: `
			Read the TASKS.md file and do the next task.
			After each task, update the TASKS.md file to reflect the changes.

			Do not leave TODOs, placeholders, etc. Fill in all the details.
			If you cannot continue, create a new task in the TASKS.md file.
			`,
		})
		tasks, err := os.ReadFile(filepath.Join(workingDirectory, "TASKS.md"))
		if err != nil {
			fmt.Printf("Error reading TASKS.md: %s", err)
			os.Exit(1)
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
			handleChatCompletion(MODEL_SMALL, openai.ChatCompletionMessage{
				Role: "user",
				Content: fmt.Sprintf(`
					Create tasks in the TASKS.md file to implement the missing functionality based on the TODOs, placeholders, etc. in the following git diff:

					git diff:
					%s
					`, string(diff)),
			})
			log.Printf("There are pending todos, continuing")
			continue
		}
		// Finished all tasks, no pending todos
		log.Printf("No more tasks, no pending todos, generating wiki")
		GenWiki()
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
