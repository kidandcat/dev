# tools.go

This file defines the tools available to the agent for system and codebase manipulation.

## Functions

- GetToolsfunc: Returns a list of all available tools.
- ToolCallfunc: Executes a specified tool based on type and parameters.
- getMultipleStringsfunc: Utility to extract a slice of strings from a map.
- getStringfunc: Utility to extract a string from a map.
- getMapfunc: Utility to extract a map from data.
- getIntfunc: Utility to extract an integer value from a map.

## Available Tools

The agent leverages tools for file operations, code editing, and web actions, defined here for extensibility and centralized management.
