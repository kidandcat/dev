# Project High Level Documentation

## Overview

This project is an autonomous agent system implemented in Go, designed to interact with the file system, codebase, and the web. The agent can read, modify, and generate documentation, respond to tasks, and update the project wiki. Its architecture leverages modular tooling and a conversational interface with the Gemini API to iteratively accomplish goals set via markdown checklists. The project provides an interactive way for the agent to self-manage tasks, improve its environment, and maintain up-to-date documentation.

## Main Components

### 1. Core Agent Logic (`agent.go`)
- Implements the decision-making loop and communications with the Gemini API.
- Coordinates task execution and tool selection.
- Maintains the context of operations using conversation history, model choice, and temperature settings for the API's response variability.

### 2. Code Manipulation (`code.go`)
- Provides functionalities to read, lint, and modify Go code.
- Enables the agent to make changes in the codebase, automatically handle imports, and ensure code quality by linting.

### 3. File System Interaction (`file.go`)
- Empowers the agent to list directory contents, read/write files, and manage directories.
- Handles .gitignore patterns and encapsulates safe file operations.

### 4. Task Management & Main Loop (`main.go`)
- Initializes the entire system and runs the agent's main execution loop.
- Monitors for pending tasks and delegates them for completion.

### 5. Tooling Abstraction (`tools.go`)
- Defines the tools available to the agent, such as file and code manipulation, and web search.
- Handles tool calls and input parsing to interface with the agent.

### 6. Web Interaction (`web.go`)
- Allows the agent to scrape web pages and perform web searches for supplementary information.

### 7. Wiki Generation (`wiki.go`)
- Automates the production and maintenance of markdown documentation in the wiki.
- Ensures documentation is kept synchronized with code changes.

## Wiki System
All components are documented in the `wiki` folder as markdown files named after their Go source files. These documentation files contain breakdowns of major types and functions, with clear parameter and return value listing, plus high-level module purposes.

## Summary
The project forms a feedback loop in which an agent autonomously manages its own file system, codebase, and documentation. Tasks are tracked in a markdown checklist, and updates to the code and docs are made as tasks are performed. The system is extensible and adapts to new requirements through code and documentation regeneration.
