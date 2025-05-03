# agent.go

This file defines the agent's core logic and interaction with the environment. It handles the agent's decision-making process, tool selection, and communication with the Gemini API.

## Constants

-   `MODEL_BIG`: The name of the larger Gemini model.
-   `MODEL_SMALL`: The name of the smaller Gemini model.

## Variables

-   `messages`: A slice of `genai.Content` representing the conversation history.
-   `workingDirectory`: The directory where the agent operates.
-   `temperature`: The temperature used for the Gemini API.

## Functions

-   `handleChatCompletionfunc`: Handles chat completions using the Gemini API.
-   `handleToolCallfunc`: Handles tool calls from the Gemini API.
-   `YesNoQuestionfunc`: Asks a yes/no question and returns the answer.

## Core Logic

The agent uses a combination of the Gemini API and available tools to accomplish tasks. It maintains a conversation history to provide context for its decisions. The `temperature` variable controls the randomness of the agent's responses.
