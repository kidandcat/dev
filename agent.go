package main

import (
	"context"
	"fmt"
	"time"

	"github.com/openai/openai-go/shared"
	openai "github.com/sashabaranov/go-openai"
)

const (
	MODEL_GPT41 = shared.ChatModel("gpt-4.1-2025-04-14")
	MODEL_NANO  = shared.ChatModel("gpt-4.1-nano-2025-04-14")
)

var messages []openai.ChatCompletionMessage

func handleChatCompletion(model string, msg openai.ChatCompletionMessage, viewModel *Model) {
	messages = append(messages, msg)

	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role: "system",
					Content: fmt.Sprintf(`
					You are an autonomous programmer agent, which can write code, fix bugs, and implement features.
					You have tools to analyze the local codebase, search the web, and more.
					Do not ask follow up questions.

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
		viewModel.AppendError(fmt.Errorf("Error creating chat completion for (%#v): %v", msg, err))
		return
	}

	messages = append(messages, response.Choices[0].Message)
	if response.Choices[0].Message.Content != "" {
		viewModel.AppendAssistant(response.Choices[0].Message.Content)
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall == response.Choices[0].Message.ToolCalls[len(response.Choices[0].Message.ToolCalls)-1] {
			handleChatCompletion(model, handleToolCall(toolCall, viewModel), viewModel)
			return
		}
		messages = append(messages, handleToolCall(toolCall, viewModel))
	}
}

func handleToolCall(toolCall openai.ToolCall, viewModel *Model) openai.ChatCompletionMessage {
	res := ToolCall(toolCall, viewModel)
	viewModel.AppendInfo(fmt.Sprintf("%s(%s) -> %s", toolCall.Function.Name, toolCall.Function.Arguments, res))
	return openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleTool,
		Content:    res,
		Name:       toolCall.Function.Name,
		ToolCallID: toolCall.ID,
	}
}
