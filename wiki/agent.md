# agent.go

This file defines the agent's core logic and interaction with the environment. It handles the agent's decision-making process, tool selection, and communication with the Gemini API.

## Constants

-   `MODEL_BIG`: The name of the larger Gemini model.
-   `MODEL_SMALL`: The name of the smaller Gemini model.

## Variables

-   `messages`: A slice of `genai.Content` representing the conversation history. This slice stores the messages exchanged between the agent and the language model, providing context for subsequent interactions.
-   `workingDirectory`: The directory where the agent operates. This variable specifies the root directory for file system operations.
-   `temperature`: The temperature used for the Gemini API. This value controls the randomness of the language model's responses. Higher values result in more random responses, while lower values result in more deterministic responses.

## Functions

-   `handleChatCompletion`: Handles chat completions using the Gemini API.
    -   **Parameters:**
        -   `model` (string): The name of the language model to use.
        -   `msg` (*genai.Content): The message to send to the language model.
    -   **Return Value:**
        -   (string): The language model's response.
    -   **Description:** This function sends a message to the language model and returns the model's response. It also handles tool calls, which allow the agent to interact with the environment.
-   `handleToolCall`: Handles tool calls from the Gemini API.
    -   **Parameters:**
        -   `toolCall` (*genai.FunctionCall): The tool call to handle.
    -   **Return Value:**
        -   (*genai.Part): The result of the tool call.
    -   **Description:** This function executes a tool call and returns the result. It uses the `ToolCall` function in `tools.go` to execute the tool.
-   `YesNoQuestion`: Asks a yes/no question and returns the answer.
    -   **Parameters:**
        -   `question` (string): The question to ask.
    -   **Return Value:**
        -   (bool): True if the answer is yes, false if the answer is no.
    -   **Description:** This function asks a yes/no question to the language model and returns the answer. It enforces a "yes" or "no" answer through tool calls.

## Core Logic

The agent uses a combination of the Gemini API and available tools to accomplish tasks. It maintains a conversation history to provide context for its decisions. The `temperature` variable controls the randomness of the agent's responses.
