# agent.go

This file defines the agent's core logic and interaction with the Gemini API.

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