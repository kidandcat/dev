package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/genai"
)

const (
	// MODEL_BIG   = shared.ChatModel("gpt-4.1-2025-04-14")
	// MODEL_SMALL = shared.ChatModel("gpt-4.1-nano-2025-04-14")

	MODEL_BIG   = "gemini-2.5-flash-preview-04-17"
	MODEL_SMALL = "gemini-2.0-flash-lite"
)

var messages []*genai.Content
var workingDirectory string
var temperature float32 = 0.7

func handleChatCompletion(model string, msg *genai.Content) string {
	messages = append(messages, msg)

	response, err := client.Models.GenerateContent(
		context.Background(),
		model,
		messages,
		&genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(fmt.Sprintf(`
					You are an autonomous, unsupervised agent that can write Go code, fix bugs, and implement features.
					You have tools to analyze the local codebase, search the web, and more.
					
					Date and time: %s
					`, time.Now().Format(time.RFC3339)), genai.RoleUser),
			Temperature: &temperature,
			Tools:       GetTools(),
		},
	)
	if err != nil {
		var contextLength int
		for _, message := range messages {
			contextLength += len(message.Parts[0].Text)
		}
		log.Printf("Context length: %d", contextLength)
		panic(err)
	}

	messages = append(messages, response.Candidates[0].Content)
	content := response.Candidates[0].Content.Parts[0].Text
	log.Printf("Assistant: %s", content)

	toolCalls := response.FunctionCalls()
	parts := []*genai.Part{}
	for _, toolCall := range toolCalls {
		if toolCall.Name == "continue" {
			return "Continue"
		}
		log.Printf("Tool call: %s", toolCall.Name)
		parts = append(parts, handleToolCall(toolCall))
	}
	if len(parts) > 0 {
		return handleChatCompletion(model, genai.NewContentFromParts(parts, genai.RoleUser))
	}

	return messages[len(messages)-1].Parts[0].Text
}

func handleToolCall(toolCall *genai.FunctionCall) *genai.Part {
	res := ToolCall(toolCall)
	return genai.NewPartFromFunctionResponse(toolCall.Name, res)
}

func YesNoQuestion(question string) bool {
	response, err := client.Models.GenerateContent(
		context.Background(),
		MODEL_SMALL,
		[]*genai.Content{
			genai.NewContentFromText(question, genai.RoleUser),
		},
		&genai.GenerateContentConfig{
			Temperature:       &temperature,
			SystemInstruction: genai.NewContentFromText(`You must answer the question with a "yes" or "no" tool call.`, genai.RoleUser),
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						{
							Name:        "yes",
							Description: "Answer affirmatively",
						},
						{
							Name:        "no",
							Description: "Answer negatively",
						},
					},
				},
			},
		},
	)
	if err != nil {
		return false
	}
	toolCalls := response.FunctionCalls()
	if len(toolCalls) == 0 {
		return false
	}
	return toolCalls[0].Name == "yes"
}
