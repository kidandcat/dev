package main

import "google.golang.org/genai"

func GenWiki() {
	handleChatCompletion(MODEL, &genai.Content{
		Role: genai.RoleUser,
		Parts: []*genai.Part{
			genai.NewPartFromText(`
			1. Analyze the code in the current directory and generate high level documentation for the code.
			2. Analyze the existing documentation in the wiki folder.
			3. Write the documentation in markdown files in the wiki folder.

			Do not leave TODOs, placeholders, etc. Fill in all the details.
			If you cannot continue, create a new task in the TASKS.md file.
			`),
		},
	})
}
