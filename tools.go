package main

import (
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func GetTools() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "visit_web_page",
				Description: "Visit a web page and get the source HTML",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"url": {
							Type:        jsonschema.String,
							Description: "The url of the web page to visit",
						},
						"headers": {
							Type:        jsonschema.Object,
							Description: "The headers to send to the web page (optional)",
						},
						"cookies": {
							Type:        jsonschema.Object,
							Description: "The cookies to send to the web page (optional)",
						},
					},
					Required: []string{"url"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "web_page_search",
				Description: "Search a web page for a query",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"query": {
							Type:        jsonschema.String,
							Description: "The query to search for",
						},
					},
					Required: []string{"query"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "list_directory",
				Description: "List the files in a directory",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to list the files in, relative to the working directory",
						},
						"depth": {
							Type:        jsonschema.Integer,
							Description: "The depth of the subdirectories to list",
						},
					},
					Required: []string{"path", "depth"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "read_file",
				Description: "Read a file",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to read the file from, relative to the working directory",
						},
						"offset": {
							Type:        jsonschema.Integer,
							Description: "The line to start reading the file from",
						},
						"length": {
							Type:        jsonschema.Integer,
							Description: "The number of lines to read, leave blank for maximum number of lines (10000)",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "write_file",
				Description: "Write to a file",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to write the file to, relative to the working directory",
						},
						"content": {
							Type:        jsonschema.String,
							Description: "The content to write to the file",
						},
					},
					Required: []string{"path", "content"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "make_directory",
				Description: "Make a directory",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to make the directory in, relative to the working directory",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "lint_file",
				Description: "Lint a Go file to check for errors",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to lint the file from, relative to the working directory",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "search_text",
				Description: "Search for text in the working directory",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"query": {
							Type:        jsonschema.String,
							Description: "The query to search for",
						},
					},
					Required: []string{"query"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "fetch_wiki_docs",
				Description: "Fetch the documentation from the wiki folder",
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "finished",
				Description: "Finished all tasks",
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "read_code",
				Description: "Read the code of specified functions from a Go file, returning the full file with only the specified functions' bodies",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to the Go file to read",
						},
						"functions": {
							Type:        jsonschema.Array,
							Description: "List of function names to keep the full body of",
							Items: &jsonschema.Definition{
								Type: jsonschema.String,
							},
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "add_or_edit_function",
				Description: "Add a new function to a Go file or edit an existing function",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"path": {
							Type:        jsonschema.String,
							Description: "The path to the Go file to modify",
						},
						"function_name": {
							Type:        jsonschema.String,
							Description: "The name of the function to add or edit",
						},
						"function_body": {
							Type:        jsonschema.String,
							Description: "The complete function body to add or replace",
						},
					},
					Required: []string{"path", "function_name", "function_body"},
				},
			},
		},
	}
}

func ToolCall(toolCall openai.ToolCall) string {
	switch toolCall.Function.Name {
	case "visit_web_page":
		var arguments struct {
			URL     string            `json:"url"`
			Headers map[string]string `json:"headers"`
			Cookies map[string]string `json:"cookies"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling url: %s", err)
		}
		return WebSource(arguments.URL, arguments.Headers, arguments.Cookies)
	case "web_page_search":
		var arguments struct {
			Query string `json:"query"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling query: %s", err)
		}
		return WebSearch(arguments.Query)
	case "list_directory":
		var arguments struct {
			Path  string `json:"path"`
			Depth int    `json:"depth"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling path: %s", err)
		}
		return ListDirectory(arguments.Path, arguments.Depth)
	case "read_file":
		var arguments struct {
			Path   string `json:"path"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling path: %s", err)
		}
		return ReadFile(arguments.Path, arguments.Offset, arguments.Length)
	case "write_file":
		var arguments struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling path: %s", err)
		}
		return WriteFile(arguments.Path, arguments.Content)
	case "make_directory":
		var arguments struct {
			Path string `json:"path"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling path: %s", err)
		}
		return MkDir(arguments.Path)
	case "fetch_wiki_docs":
		return FetchWikiDocs()
	case "lint_file":
		var arguments struct {
			Path string `json:"path"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling path: %s", err)
		}
		return Lint(arguments.Path)
	case "search_text":
		var arguments struct {
			Query string `json:"query"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling query: %s", err)
		}
		return SearchText(arguments.Query)
	case "read_code":
		var arguments struct {
			Path      string   `json:"path"`
			Functions []string `json:"functions"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling arguments: %s", err)
		}
		return ReadCode(arguments.Path, arguments.Functions...)
	case "add_or_edit_function":
		var arguments struct {
			Path         string `json:"path"`
			FunctionName string `json:"function_name"`
			FunctionBody string `json:"function_body"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling arguments: %s", err)
		}
		return AddOrEditFunction(arguments.Path, arguments.FunctionName, arguments.FunctionBody)
	}
	return fmt.Sprintf("Unknown tool call: %s", toolCall.Function.Name)
}
