package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
)

var workingDirectory string

type Model struct {
	messages []string
	aiModel  string
}

func NewModel() *Model {
	return &Model{
		messages: []string{},
		aiModel:  MODEL_NANO,
	}
}

func (m *Model) Start() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("You: ")
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)
		m = m.AppendUser(content)
		handleChatCompletion(m.aiModel, openai.ChatCompletionMessage{
			Role:    "user",
			Content: content,
		}, m)
	}
}

func (m *Model) AppendUser(msg string) *Model {
	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	m.messages = append(m.messages, userStyle.Render("You: ")+msg)
	return m
}

func (m *Model) AppendInfo(msg string) *Model {
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	m.messages = append(m.messages, infoStyle.Render("Info: ")+msg)
	return m
}

func (m *Model) AppendAssistant(msg string) *Model {
	assistantStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	m.messages = append(m.messages, assistantStyle.Render("Assistant: ")+msg)
	return m
}

func (m *Model) AppendError(err error) *Model {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")). // Red color
		Bold(true)

	msg := errorStyle.Render(fmt.Sprintf("[ERROR] %v", err))
	m.messages = append(m.messages, msg)
	return m
}
