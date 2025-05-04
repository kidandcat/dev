package main

import (
	"fmt"
	"strconv"

	"google.golang.org/genai"
)

func GetTools() []*genai.Tool {
	return []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "visit_web_page",
					Description: "Visit a web page and get the source HTML",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"url":     {Type: genai.TypeString, Description: "The url of the web page to visit"},
							"headers": {Type: genai.TypeObject, Description: "The headers to send to the web page (optional)"},
							"cookies": {Type: genai.TypeObject, Description: "The cookies to send to the web page (optional)"},
						},
						Required: []string{"url"},
					},
				},
				{
					Name:        "web_page_search",
					Description: "Search a web page for a query",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {Type: genai.TypeString, Description: "The query to search for"},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "list_directory",
					Description: "List the files in a directory",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path":  {Type: genai.TypeString, Description: "The path to list the files in, relative to the working directory"},
							"depth": {Type: genai.TypeInteger, Description: "The depth of the subdirectories to list"},
						},
						Required: []string{"path", "depth"},
					},
				},
				{
					Name:        "read_file",
					Description: "Read a file",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path":   {Type: genai.TypeString, Description: "The path to read the file from, relative to the working directory"},
							"offset": {Type: genai.TypeInteger, Description: "The line to start reading the file from"},
							"length": {Type: genai.TypeInteger, Description: "The number of lines to read"},
						},
						Required: []string{"path"},
					},
				},
				{
					Name:        "write_file",
					Description: "Write to a file",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path":    {Type: genai.TypeString, Description: "The path to write the file to, relative to the working directory"},
							"content": {Type: genai.TypeString, Description: "The content to write to the file"},
						},
						Required: []string{"path", "content"},
					},
				},
				{
					Name:        "make_directory",
					Description: "Create a directory",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path": {Type: genai.TypeString, Description: "The path to make the directory in, relative to the working directory"},
						},
						Required: []string{"path"},
					},
				},
				{
					Name:        "lint_file",
					Description: "Lint a Go file to check for errors",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path": {Type: genai.TypeString, Description: "The path to lint the file from, relative to the working directory"},
						},
						Required: []string{"path"},
					},
				},
				{
					Name:        "search_text",
					Description: "Search for text in the working directory",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"query": {Type: genai.TypeString, Description: "The query to search for"},
						},
						Required: []string{"query"},
					},
				},
				{
					Name:        "fetch_wiki_docs",
					Description: "Fetch the documentation from the wiki folder",
					Parameters: &genai.Schema{
						Type:       genai.TypeObject,
						Properties: map[string]*genai.Schema{},
						Required:   []string{},
					},
				},
				{
					Name:        "continue",
					Description: "Task finished, continue with the next task",
					Parameters: &genai.Schema{
						Type:       genai.TypeObject,
						Properties: map[string]*genai.Schema{},
						Required:   []string{},
					},
				},
				{
					Name:        "read_code",
					Description: "Read the code of specified functions from a Go file, returning the full file with only the specified functions' bodies",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path": {Type: genai.TypeString, Description: "The path to the Go file to read, relative to the working directory"},
							"functions": {
								Type:        genai.TypeArray,
								Description: "The list of function names to read the full body of.",
								Items: &genai.Schema{
									Type: genai.TypeString,
								},
							},
						},
						Required: []string{"path"},
					},
				},
				{
					Name:        "add_or_edit_function",
					Description: "Add a new function to a Go file or edit an existing function",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"path":          {Type: genai.TypeString, Description: "The path to the Go file to modify, relative to the working directory"},
							"function_name": {Type: genai.TypeString, Description: "The name of the function to add or edit"},
							"function_body": {Type: genai.TypeString, Description: "The code that will go inside the function body"},
						},
						Required: []string{"path", "function_name", "function_body"},
					},
				},
			},
		},
	}
}

func ToolCall(toolCall *genai.FunctionCall) map[string]any {
	switch toolCall.Name {
	case "visit_web_page":
		return WebSource(getString(toolCall.Args["url"]), getMap(toolCall.Args["headers"]), getMap(toolCall.Args["cookies"]))
	case "web_page_search":
		return WebSearch(getString(toolCall.Args["query"]))
	case "list_directory":
		return ListDirectory(getString(toolCall.Args["path"]), getInt(toolCall.Args["depth"]))
	case "read_file":
		return ReadFile(getString(toolCall.Args["path"]), getInt(toolCall.Args["offset"]), getInt(toolCall.Args["length"]))
	case "write_file":
		return WriteFile(getString(toolCall.Args["path"]), getString(toolCall.Args["content"]))
	case "make_directory":
		return MkDir(getString(toolCall.Args["path"]))
	case "fetch_wiki_docs":
		return FetchWikiDocs()
	case "lint_file":
		return Lint(getString(toolCall.Args["path"]))
	case "search_text":
		return SearchText(getString(toolCall.Args["query"]))
	case "read_code":
		return ReadCode(getString(toolCall.Args["path"]), getMultipleStrings(toolCall.Args["functions"])...)
	case "add_or_edit_function":
		return AddOrEditFunction(getString(toolCall.Args["path"]), getString(toolCall.Args["function_name"]), getString(toolCall.Args["function_body"]))
	}
	return map[string]any{
		"error": fmt.Sprintf("Unknown tool call: %s", toolCall.Name),
	}
}

func getMultipleStrings(data any) []string {
	if data == nil {
		return []string{}
	}
	switch v := data.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []any:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = getString(item)
		}
		return result
	default:
		return []string{}
	}
}

func getString(data any) string {
	if data == nil {
		return ""
	}
	switch v := data.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

func getMap(data any) map[string]string {
	if data == nil {
		return map[string]string{}
	}
	return data.(map[string]string)
}

func getInt(data any) int {
	if data == nil {
		return 0
	}
	switch v := data.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		return 0
	}
}
