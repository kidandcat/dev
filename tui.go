package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	openai "github.com/sashabaranov/go-openai"
)

const gap = "\n\n"

var workingDirectory string

type Model struct {
	viewport viewport.Model
	messages []string
	textarea textarea.Model
}

func initialModel() *Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	ta.KeyMap.InsertNewline.SetEnabled(false)

	vp := viewport.New(30, 5)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62"))

	return &Model{
		textarea: ta,
		messages: []string{},
		viewport: vp,
	}
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	m.textarea.Prompt = "┃ "

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
		if len(m.messages) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			content := m.textarea.Value()
			m = m.AppendUser(content)
			m.textarea.Reset()
			m.viewport.GotoBottom()
			go handleChatCompletion(MODEL_NANO, openai.ChatCompletionMessage{
				Role:    "user",
				Content: content,
			}, m)
		}
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *Model) AppendUser(msg string) *Model {
	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	m.messages = append(m.messages, userStyle.Render("You: ")+msg)
	// Set the border color to the assistant color
	m.viewport.Style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("2"))
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.viewport.GotoBottom()
	return m
}

func (m *Model) AppendInfo(msg string) *Model {
	infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	m.messages = append(m.messages, infoStyle.Render("Info: ")+msg)
	// Set the border color to the info color
	m.viewport.Style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("3"))
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.viewport.GotoBottom()
	return m
}

func (m *Model) AppendAssistant(msg string) *Model {
	assistantStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	m.messages = append(m.messages, assistantStyle.Render("Assistant: ")+msg)
	// Set the border color to the user color
	m.viewport.Style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("5"))
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.viewport.GotoBottom()
	return m
}

func (m *Model) AppendError(err error) *Model {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")). // Red color
		Bold(true)

	// Get the viewport width to use for wrapping
	width := m.viewport.Width
	if width == 0 {
		width = 80 // Default width if viewport width is not set
	}

	// Wrap the error message to fit the viewport width
	wrappedMsg := lipgloss.NewStyle().
		Width(width - 10). // Leave some margin
		Render(fmt.Sprintf("[ERROR] %v", err))

	msg := errorStyle.Render(wrappedMsg)
	m.messages = append(m.messages, msg)
	m.viewport.Style = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("9"))
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.viewport.GotoBottom()
	return m
}

func (m *Model) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}
