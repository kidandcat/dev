package main

import "github.com/openai/openai-go/shared"

const (
	MODEL        = shared.ChatModel("gpt-4.1-nano-2025-04-14")
	MAX_MESSAGES = 30

	MISTRAL = "mistral"
	GEMMA3  = "gemma3"
	LLAMA3  = "llama3"
	LLAMA3V = "llama3.2-vision"
)
