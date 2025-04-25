package main

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
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
							Description: "The depth to list the files in",
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
							Description: "The offset to start reading the file from",
						},
						"length": {
							Type:        jsonschema.Integer,
							Description: "The length of the file to read, 0 for all",
						},
					},
					Required: []string{"path", "offset", "length"},
				},
			},
		},
	}
}

func ToolCall(ctx context.Context, toolCall openai.ToolCall) string {
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
		source := WebSource(arguments.URL, arguments.Headers, arguments.Cookies)
		return fmt.Sprintf("Source:\n%s", source)
	case "web_page_search":
		var arguments struct {
			Query string `json:"query"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling query: %s", err)
		}
		source := WebSearch(arguments.Query)
		return fmt.Sprintf("Source:\n%+v", source)
	case "list_directory":
		var arguments struct {
			Path  string `json:"path"`
			Depth int    `json:"depth"`
		}
		err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
		if err != nil {
			return fmt.Sprintf("Error unmarshalling path: %s", err)
		}
		files := ListDirectory(arguments.Path, arguments.Depth)
		return fmt.Sprintf("Files:\n%+v", files)
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
		content := ReadFile(arguments.Path, arguments.Offset, arguments.Length)
		return fmt.Sprintf("Content:\n%s", content)
	}
	return fmt.Sprintf("Unknown tool call: %s", toolCall.Function.Name)
}
