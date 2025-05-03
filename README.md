# Autonomous Agent

This project implements an autonomous agent that can perform tasks based on instructions provided in `INPUT.md`. It uses the Gemini API to process tasks, manage files, and interact with the environment.

## How it Works

The agent operates in a loop:

1.  Reads tasks from `TASKS.md`.
2.  Processes each task using the Gemini API.
3.  Updates `TASKS.md` to reflect the completion status.
4.  If there are TODOs in the code, it creates new tasks in `TASKS.md` to address them.
5.  Once all tasks are complete and there are no TODOs, it generates a wiki.

## Usage

1.  Set the `GEMINI_API_KEY` environment variable.
2.  Create an `INPUT.md` file with a list of tasks.
3.  Run the agent: `go run main.go [working_directory]` (optional working directory).

## Files

*   `INPUT.md`: Contains the initial list of tasks.
*   `TASKS.md`: Contains the current list of tasks with completion status.
*   `main.go`: The main entry point of the application.
*   `agent.go`: Contains the agent's core logic.
*   `tools.go`: Defines the available tools for the agent.
*   `wiki.go`: Generates the project wiki.

## License

MIT