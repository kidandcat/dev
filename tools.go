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
	}
	return fmt.Sprintf("Unknown tool call: %s", toolCall.Function.Name)
}
