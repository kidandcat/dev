package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	openai "github.com/sashabaranov/go-openai"
)

func Lint(path string) string {
	path = Path(path)

	extension := filepath.Ext(path)
	dir := filepath.Dir(path)

	switch extension {
	case ".go":
		command := exec.Command("go", "vet", dir)
		output, err := command.CombinedOutput()
		if err != nil && len(output) == 0 {
			return fmt.Sprintf("Error formatting go file: %s", err)
		}
		return string(output)
	}

	return "Linting not supported for this file type"
}

func Plan(feature string, viewModel *Model) string {
	handleChatCompletion(MODEL_GPT41, openai.ChatCompletionMessage{
		Role: "user",
		Content: fmt.Sprintf(`Create a plan to develop the feature provided by the user.
		You should create a list of steps to develop the feature.
		The steps must be markdown unchecked checkboxes.
		Save the plan in a file called PLAN.md
		Feature: %s
		`, feature),
	}, viewModel, true)
	if _, err := os.Stat("PLAN.md"); os.IsNotExist(err) {
		return "Plan not created"
	}
	plan, err := os.ReadFile("PLAN.md")
	if err != nil {
		return "Plan not created"
	}
	return string(plan)
}
