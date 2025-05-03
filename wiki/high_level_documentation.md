# High Level Documentation

## Data structures and relationships

### Gitignore
- `Gitignore` is a struct defined in `file.go`.
- **Fields:**
  - `Patterns []string`: Holds a list of ignore patterns parsed from a `.gitignore` file.
- **Relationships:** Used by file and directory listing functions to determine which paths to ignore.

### WebSearchResult
- `WebSearchResult` is a struct defined in `web.go`.
- **Fields:**
  - `Position int`: The position of the result in the search results.
  - `Title string`: The title of the web result.
  - `Link string`: The hyperlink to the result.
  - `Snippet string`: A snippet/summary of the result.
- **Relationships:** Used in web search functions to collect and structure search results.

*Other data structures may exist in additional files or may be added as the project grows.*

## Methods and functions

### agent.go
- `handleChatCompletion(model, msg)`: Handles a chat completion request with the specified model and input message.
- `handleToolCall(toolCall)`: Handles a tool call and returns a chat completion message.

### code.go
- `Lint(path)`: Lints the Go file at the specified path and returns linting results.

### file.go
- `NewGitignore(path)`: Reads a `.gitignore` file and returns a `Gitignore` struct.
- `(g *Gitignore) IsIgnored(path)`: Checks if a path matches any ignore pattern.
- `ListDirectory(path, depth)`: Lists files and directories, obeying ignore patterns and depth.
- `ReadFile(path, offset, length)`: Reads a portion of a file.
- `WriteFile(path, content)`: Writes content to a file.
- `MkDir(path)`: Makes a new directory.
- `SearchText(query)`: Searches for text in files.
- `Path(path)`: Normalizes or processes a path.

### main.go
- `main()`: Entry point, dispatches chat completion and task execution logic.

### tools.go
- `GetTools()`: Returns a list of API tools.
- `ToolCall(toolCall)`: Handles calling a tool.

### web.go
- `WebSource(url, headers, cookies)`: Fetches the source HTML of a webpage.
- `WebSearch(query)`: Performs a web search and structures the results.

*Test functions and other helpers exist in `*_test.go` files.*

## Files and folders

### Top-level files
- `INPUT.md`: Input file for new tasks.
- `TASKS.md`: The task checklist for the project.
- `go.mod`, `go.sum`: Go module and dependency files.
- `main.go`: Entry point for the application.
- `agent.go`: Handles chat completions and tool calls.
- `code.go`: Contains code linting utilities.
- `file.go`: File handling, directory listing, and file operations.
- `tools.go`: API tool definitions and handlers.
- `web.go`: Web source fetching and search logic.

### Test files
- `code_test.go`, `file_test.go`: Unit tests for their corresponding implementations.

### Folders
- `demo/`: Example or demonstration files. Contains files like `main.go`, `names.db`, and `public/` directory for demo assets.
- `testdata/`: Contains test fixtures and supporting files for tests (e.g., `empty`, `gitignore`, `gitignore_complex`, `hidden`, `nested`, `single`).
- `wiki/`: Project documentation. Includes `high_level_documentation.md`.

*Other files and folders may be added as needed.*
