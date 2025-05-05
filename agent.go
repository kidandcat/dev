package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	// MODEL = "google/gemini-2.0-flash-001"
	// MODEL = "x-ai/grok-3-mini-beta"
	MODEL = "anthropic/claude-3.5-haiku"
)

var messages []openai.ChatCompletionMessage
var workingDirectory string

func handleChatCompletion(model string, msg openai.ChatCompletionMessage) string {
	messages = append(messages, msg)

	if len(messages) > 10 {
		messages = messages[len(messages)-10:]
	}

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf(`
					You are an autonomous, unsupervised agent that can write Go code, fix bugs, and implement features.
					You have tools to analyze the local codebase, search the web, and more.
					
					Date and time: %s
					`, time.Now().Format(time.RFC3339)),
				},
			}, messages...),
			Tools: GetTools(),
		},
	)
	if err != nil {
		var contextLength int
		for _, message := range messages {
			contextLength += len(message.Content)
		}
		log.Printf("Context length: %d", contextLength)
		panic(err)
	}

	if len(response.Choices) == 0 || (response.Choices[0].Message.Content == "" && response.Choices[0].Message.ToolCalls == nil) {
		panic(fmt.Sprintf("No response from assistant: %+v", response))
	}

	messages = append(messages, response.Choices[0].Message)
	if response.Choices[0].Message.Content != "" {
		log.Printf("Assistant: %s", response.Choices[0].Message.Content)
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == "finished" {
			return "Finished all tasks"
		}
		if toolCall == response.Choices[0].Message.ToolCalls[len(response.Choices[0].Message.ToolCalls)-1] {
			return handleChatCompletion(model, handleToolCall(toolCall))
		}
		messages = append(messages, handleToolCall(toolCall))
	}

	return messages[len(messages)-1].Content
}

func handleToolCall(toolCall openai.ToolCall) openai.ChatCompletionMessage {
	res := ToolCall(toolCall)
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    res,
		ToolCallID: toolCall.ID,
	}
}

func YesNoQuestion(question string) bool {
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: MODEL,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "system",
					Content: `You must answer the question with a "yes" or "no" tool call.`,
				},
				{
					Role:    "user",
					Content: question,
				},
			},
			Stream:      false,
			Temperature: 0.7,
			Tools: []openai.Tool{
				{
					Type: openai.ToolTypeFunction,
					Function: &openai.FunctionDefinition{
						Name:        "yes",
						Description: "Answer affirmatively",
					},
				},
				{
					Type: openai.ToolTypeFunction,
					Function: &openai.FunctionDefinition{
						Name:        "no",
						Description: "Answer negatively",
					},
				},
			},
		},
	)
	if err != nil {
		return false
	}
	if response.Choices[0].Message.ToolCalls[0].Function.Name == "yes" {
		return true
	}
	if response.Choices[0].Message.ToolCalls[0].Function.Name == "no" {
		return false
	}
	if strings.Contains(response.Choices[0].Message.Content, "yes") {
		return true
	}
	return false
}
