package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/openai/openai-go/shared"
	openai "github.com/sashabaranov/go-openai"
)

const (
	MODEL_GPT41 = shared.ChatModel("gpt-4.1-2025-04-14")
	MODEL_NANO  = shared.ChatModel("gpt-4.1-nano-2025-04-14")
)

var messages []openai.ChatCompletionMessage
var workingDirectory string

func handleChatCompletion(model string, msg openai.ChatCompletionMessage) string {
	messages = append(messages, msg)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: fmt.Sprintf(`
					You are an autonomous, unsupervised agent that can write Go code, fix bugs, and implement features.
					There is no user able to provide feedback. You are the only one in the conversation.
					You have tools to analyze the local codebase, search the web, and more.
					
					Date and time: %s
					`, time.Now().Format(time.RFC3339)),
				},
			}, messages...),
			Stream:      false,
			Temperature: 0.7,
			Tools:       GetTools(),
		},
	)
	if err != nil {
		log.Printf("Error creating chat completion for (%#v): %v", msg, err)
		return fmt.Sprintf("Error creating chat completion for (%#v): %v", msg, err)
	}

	messages = append(messages, response.Choices[0].Message)
	if response.Choices[0].Message.Content != "" {
		log.Printf("Assistant: %s", response.Choices[0].Message.Content)
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == "continue" {
			return "Continue"
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
	log.Printf("%s(%s)", toolCall.Function.Name, toolCall.Function.Arguments)
	if res == "" {
		panic(fmt.Sprintf("Tool call %s(%s) returned empty string", toolCall.Function.Name, toolCall.Function.Arguments))
	}
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    res,
		Name:       toolCall.Function.Name,
		ToolCallID: toolCall.ID,
	}
}

func YesNoQuestion(question string) bool {
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: MODEL_NANO,
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
