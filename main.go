package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	openai "github.com/sashabaranov/go-openai"
)

var client *openai.Client

func main() {
	client = openai.NewClient("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")

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
