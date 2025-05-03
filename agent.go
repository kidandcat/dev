package main

import (
	"context"
	"fmt"
	"log"
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

func handleChatCompletion(model string, msg openai.ChatCompletionMessage) {
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
		return
	}

	messages = append(messages, response.Choices[0].Message)
	if response.Choices[0].Message.Content != "" {
		log.Printf("Assistant: %s", response.Choices[0].Message.Content)
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall.Function.Name == "continue" {
			log.Println("Continue")
			return
		}
		if toolCall == response.Choices[0].Message.ToolCalls[len(response.Choices[0].Message.ToolCalls)-1] {
			handleChatCompletion(model, handleToolCall(toolCall))
			return
		}
		messages = append(messages, handleToolCall(toolCall))
	}
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
