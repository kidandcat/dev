package main

import (
	"github.com/sashabaranov/go-openai"
)

func GenWiki() {
	handleChatCompletion(MODEL, openai.ChatCompletionMessage{
		Role: "user",
		Content: `
		1. Analyze the code in the current directory and generate high level documentation for the code.
		2. Analyze the existing documentation in the wiki folder.
		3. Write the documentation in markdown files in the wiki folder.

		Do not leave TODOs, placeholders, etc. Fill in all the details.
		If you cannot continue, create a new task in the TASKS.md file.
		`,
	})
}
