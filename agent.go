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

					After you have written code, always use the lint_file tool to check if the code is correct.
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
			handleChatCompletion(ctx, handleToolCall(ctx, toolCall, model), model)
			return
		}
		messages = append(messages, handleToolCall(ctx, toolCall, model))
	}
}

func handleToolCall(ctx context.Context, toolCall openai.ToolCall, model *Model) openai.ChatCompletionMessage {
	res := ToolCall(ctx, toolCall)
	model.AppendInfo(fmt.Sprintf("%s(%s) -> %s", toolCall.Function.Name, toolCall.Function.Arguments, res))
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    res,
		Name:       toolCall.Function.Name,
		ToolCallID: toolCall.ID,
	}
}
