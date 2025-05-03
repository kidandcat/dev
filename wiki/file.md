# file.go

This file contains functions for interacting with the file system.

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