# file.go

This file contains functions for interacting with the file system. These functions allow the agent to read, write, and manipulate files and directories.

## Types

-   Gitignore: Represents a .gitignore file with a slice of patterns for exclusions.

## Functions

-   NewGitignorefunc: Creates and loads ignore rules from a .gitignore file.
-   IsIgnoredfunc: Checks if a file/path is excluded by .gitignore rules.
-   ListDirectoryfunc: Lists files (recursively by depth) in a directory, respecting .gitignore.
-   ReadFilefunc: Reads file content by offset and length, preventing large or Go file reads directly.
-   WriteFilefunc: Writes text to a file by path.
-   MkDirfunc: Creates new directories.
-   FetchWikiDocsfunc: Collects documentation from the wiki folder.
-   SearchTextfunc: Recursively searches for a query in files.
-   Pathfunc: Resolves a path string relative to the working directory.

## File System Interaction

These functions are foundational for the agent's file navigation and management, ensuring it can read, write, and organize resources safely and according to ignore rules.
