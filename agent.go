package main

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var messages []openai.ChatCompletionMessage

func handleChatCompletion(ctx context.Context, msg openai.ChatCompletionMessage, model *Model) {
	messages = append(messages, msg)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: MODEL,
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: fmt.Sprintf(`
					You are an autonomous programmer agent, which can write code, fix bugs, and implement features.
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
		model.AppendError(fmt.Errorf("Error creating chat completion for (%#v): %v", msg, err))
		return
	}

	messages = append(messages, response.Choices[0].Message)
	if response.Choices[0].Message.Content != "" {
		model.AppendAssistant(response.Choices[0].Message.Content)
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall == response.Choices[0].Message.ToolCalls[len(response.Choices[0].Message.ToolCalls)-1] {
			handleChatCompletion(ctx, handleToolCall(ctx, toolCall), model)
			return
		}
		messages = append(messages, handleToolCall(ctx, toolCall))
	}
}

func handleToolCall(ctx context.Context, toolCall openai.ToolCall) openai.ChatCompletionMessage {
	res := ToolCall(ctx, toolCall)
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    res,
		Name:       toolCall.Function.Name,
		ToolCallID: toolCall.ID,
	}
}
