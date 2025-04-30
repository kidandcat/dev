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
					You are an autonomous agent that can write code, fix bugs, and implement features.
					You have tools to analyze the local codebase, search the web, and more.

					Check if a file TASKS.md exists in the working directory.
					If it does, read it and do the next task.
					
					If it doesn't, open a file called INPUT.md and read the content.
					Create a list of tasks to complete based on the content of the INPUT.md file.
					Write the markdown taskslist to a file called TASKS.md.

					Then, iterate over the tasks in TASKS.md and do them one by one.
					After each task, update the TASKS.md file to reflect the changes.


					Use the tool "finished" to finish the program when all tasks are completed.

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
		logger.Printf("Error creating chat completion for (%#v): %v", msg, err)
		return
	}

	messages = append(messages, response.Choices[0].Message)
	if response.Choices[0].Message.Content != "" {
		logger.Printf("Assistant: %s", response.Choices[0].Message.Content)
	}

	for _, toolCall := range response.Choices[0].Message.ToolCalls {
		if toolCall == response.Choices[0].Message.ToolCalls[len(response.Choices[0].Message.ToolCalls)-1] {
			handleChatCompletion(model, handleToolCall(toolCall))
			return
		}
		messages = append(messages, handleToolCall(toolCall))
	}
}

func handleToolCall(toolCall openai.ToolCall) openai.ChatCompletionMessage {
	res := ToolCall(toolCall)
	logger.Printf("%s(%s)", toolCall.Function.Name, toolCall.Function.Arguments)
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
