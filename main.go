package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
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
		os.Create(filepath.Join(workingDirectory, "INPUT.md"))
		os.Exit(1)
	}

	if _, err := os.Stat(filepath.Join(workingDirectory, "TASKS.md")); os.IsNotExist(err) {
		os.Create(filepath.Join(workingDirectory, "TASKS.md"))
	}

	handleChatCompletion(MODEL, &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			genai.NewPartFromText(`
			Open a file called INPUT.md and read the content.
			Process the content of the INPUT.md file into independent, small tasks and add them to the markdown checklist in the TASKS.md file.

			If the file TASKS.md does not exist, create it.
			`),
		},
	})

	for {
		response := handleChatCompletion(MODEL, &genai.Content{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				genai.NewPartFromText(`
				Read the TASKS.md file and do the next task.
				After each task, update the TASKS.md file to reflect the changes.

				Do not leave TODOs, placeholders, etc. Fill in all the details.
				If you cannot continue, create a new task in the TASKS.md file.
				`),
			},
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
			handleChatCompletion(MODEL, &genai.Content{
				Role: genai.RoleUser,
				Parts: []*genai.Part{
					genai.NewPartFromText(fmt.Sprintf(`
					Create tasks in the TASKS.md file to implement the missing functionality based on the TODOs, placeholders, etc. in the following git diff:

					git diff:
					%s
					`, string(diff))),
				},
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
		// Erase the TASKS.md file
		if err := os.WriteFile(filepath.Join(workingDirectory, "TASKS.md"), []byte{}, 0644); err != nil {
			fmt.Printf("Error erasing TASKS.md: %s", err)
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
