# file.go

This file contains functions for interacting with the file system. These functions allow the agent to read, write, and manipulate files and directories.

## Types

-   `Gitignore`: Represents a .gitignore file.
    -   **Fields:**
        -   `Patterns` ([]string): A list of ignore patterns.

## Functions

-   `NewGitignore`: Creates a new `Gitignore` object.
    -   **Parameters:**
        -   `path` (string): The path to the .gitignore file.
    -   **Return Value:**
        -   (*Gitignore): A pointer to the new `Gitignore` object.
    -   **Description:** This function creates a new `Gitignore` object by reading the ignore patterns from the specified .gitignore file.
-   `IsIgnored`: Checks if a path is ignored by a .gitignore file.
    -   **Parameters:**
        -   `path` (string): The path to check.
    -   **Return Value:**
        -   (bool): True if the path is ignored, false otherwise.
    -   **Description:** This function checks if the specified path is ignored based on the patterns in the `Gitignore` object.
-   `ListDirectory`: Lists the files in a directory.
    -   **Parameters:**
        -   `path` (string): The path to the directory to list.
        -   `depth` (int): The depth of subdirectories to list.
    -   **Return Value:**
        -   (map[string]any): A map containing the list of files. The map will contain a "files" key with a list of file names as a string array.
    -   **Description:** This function lists the files in the specified directory, recursively up to the given depth. It also respects .gitignore files.
-   `ReadFile`: Reads a file.
    -   **Parameters:**
        -   `path` (string): The path to the file to read.
        -   `offset` (int): The line to start reading the file from.
        -   `length` (int): The number of lines to read.
    -   **Return Value:**
        -   (map[string]any): A map containing the content of the file. The map will contain a "content" key with the file content as a string.
    -   **Description:** This function reads the specified file and returns its content. It limits the number of lines that can be read to 1000 and prevents reading Go files directly (suggesting the use of `code.ReadCode` instead).
-   `WriteFile`: Writes to a file.
    -   **Parameters:**
        -   `path` (string): The path to the file to write.
        -   `content` (string): The content to write to the file.
    -   **Return Value:**
        -   (map[string]any): A map containing the results of the operation. The map will contain a "path" key with the path to the file, a "content" key with the content of the file, and a "lint" key with the linting results.
# file.go

This file contains functions for interacting with the file system. These functions allow the agent to read, write, and manipulate files and directories.

## Types

-   `Gitignore`: Represents a .gitignore file.
    -   **Fields:**
        -   `Patterns` ([]string): A list of ignore patterns.

## Functions

-   `NewGitignore`: Creates a new `Gitignore` object.
    -   **Parameters:**
        -   `path` (string): The path to the .gitignore file.
    -   **Return Value:**
        -   (*Gitignore): A pointer to the new `Gitignore` object.
    -   **Description:** This function creates a new `Gitignore` object by reading the ignore patterns from the specified .gitignore file.
-   `IsIgnored`: Checks if a path is ignored by a .gitignore file.
    -   **Parameters:**
        -   `path` (string): The path to check.
    -   **Return Value:**
        -   (bool): True if the path is ignored, false otherwise.
    -   **Description:** This function checks if the specified path is ignored based on the patterns in the `Gitignore` object.
-   `ListDirectory`: Lists the files in a directory.
    -   **Parameters:**
        -   `path` (string): The path to the directory to list.
        -   `depth` (int): The depth of subdirectories to list.
    -   **Return Value:**
        -   (map[string]any): A map containing the list of files. The map will contain a "files" key with a list of file names as a string array.
    -   **Description:** This function lists the files in the specified directory, recursively up to the given depth. It also respects .gitignore files.
-   `ReadFile`: Reads a file.
    -   **Parameters:**
        -   `path` (string): The path to the file to read.
        -   `offset` (int): The line to start reading the file from.
        -   `length` (int): The number of lines to read.
    -   **Return Value:**
        -   (map[string]any): A map containing the content of the file. The map will contain a "content" key with the file content as a string.
    -   **Description:** This function reads the specified file and returns its content. It limits the number of lines that can be read to 1000 and prevents reading Go files directly (suggesting the use of `code.ReadCode` instead).
-   `WriteFile`: Writes to a file.
    -   **Parameters:**
        -   `path` (string): The path to the file to write.
        -   `content` (string): The content to write to the file.
    -   **Return Value:**
        -   (map[string]any): A map containing the results of the operation. The map will contain a "path" key with the path to the file, a "content" key with the content of the file, and a "lint" key with the linting results.
# file.go

This file contains functions for interacting with the file system. These functions allow the agent to read, write, and manipulate files and directories.

## Types

-   `Gitignore`: Represents a .gitignore file.

## Functions

-   `NewGitignorefunc`: Creates a new `Gitignore` object.
-   `IsIgnoredfunc`: Checks if a path is ignored by a .gitignore file.
-   `ListDirectoryfunc`: Lists the files in a directory.
-   `ReadFilefunc`: Reads a file.
-   `WriteFilefunc`: Writes to a file.
-   `MkDirfunc`: Creates a directory.
-   `FetchWikiDocsfunc`: Fetches documentation from the wiki folder.
-   `SearchTextfunc`: Searches for text in the working directory.
-   `searchTextRecursivefunc`: Recursively searches for text in a directory.
-   `Pathfunc`: Resolves a path relative to the working directory.

## File System Interaction

The functions in this file provide the agent with the ability to navigate the file system, read file contents, write new files, and create directories. This is essential for the agent's ability to modify the codebase and manage its own files.
