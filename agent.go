package main

import (
	"context"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	// MODEL = "anthropic/claude-3.7-sonnet"
	MODEL = "openai/gpt-4o-mini"
)

var messages []openai.ChatCompletionMessage
var workingDirectory string

func handleChatCompletion(msg openai.ChatCompletionMessage) string {
	if len(messages) > 10 {
		messages = messages[len(messages)-10:]
	}

	pendingMessages := []openai.ChatCompletionMessage{
		msg,
	}

	for len(pendingMessages) > 0 {
		messages = append(messages, pendingMessages...)
		pendingMessages = nil

		response, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:    MODEL,
				Messages: messages,
				Tools:    GetTools(),
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
			log.Printf("No response from assistant: %+v\n%+v\n", response, messages)
			break
		}

		messages = append(messages, response.Choices[0].Message)
		if response.Choices[0].Message.Content != "" {
			log.Printf("Assistant: %s", response.Choices[0].Message.Content)
		}

		toolCalls := response.Choices[0].Message.ToolCalls
		for _, toolCall := range toolCalls {
			if toolCall.Function.Name == "finished" {
				return "Finished all tasks"
			}
			pendingMessages = append(pendingMessages, handleToolCall(toolCall))
		}
	}

	log.Printf("Finished loop, returning last message: %s", messages[len(messages)-1].Content)
	return messages[len(messages)-1].Content
}

func handleToolCall(toolCall openai.ToolCall) openai.ChatCompletionMessage {
	log.Printf("[TOOL] %s %s", toolCall.Function.Name, toolCall.Function.Arguments)
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
		log.Printf("Error in YesNoQuestion: %s", err)
		return false
	}
	if len(response.Choices) == 0 || len(response.Choices[0].Message.ToolCalls) == 0 {
		log.Printf("No response from assistant: %+v", response)
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
