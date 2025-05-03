package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLint(t *testing.T) {
	// Define test cases
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Valid file with no errors",
			path: "./some_valid_file.go",
			want: "No errors found",
		},
		{
			name: "Non-existent path",
			path: "./non_existent.go",
			want: "Error formatting go file: open /Users/jairo/prog/non_existent.go: no such file or directory",
		},
		{
			name: "Directory with issues",
			path: "./project_with_issues",
			want: "some issues reported",
		},
		{
			name: "Clean directory",
			path: "./clean_project",
			want: "No errors found",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lint(tt.path)
			if got != tt.want {
				t.Errorf("Lint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadCode(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")

	// Write test content
	testContent := `package main

import "fmt"

type TestStruct struct {
	Field1 string
	Field2 int
}

func Function1() {
	fmt.Println("Function1 body")
}

func Function2() {
	fmt.Println("Function2 body")
}

func Function3() {
	fmt.Println("Function3 body")
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test cases
	tests := []struct {
		name      string
		functions []string
		want      string
	}{
		{
			name:      "Read all functions",
			functions: []string{"Function1", "Function2", "Function3"},
			want:      testContent,
		},
		{
			name:      "Read specific functions",
			functions: []string{"Function1", "Function3"},
			want: `package main

import "fmt"

type TestStruct struct {
	Field1 string
	Field2 int
}

func Function1() {
	fmt.Println("Function1 body")
}

func Function2() {}

func Function3() {
	fmt.Println("Function3 body")
}
`,
		},
		{
			name:      "Read no functions",
			functions: []string{},
			want: `package main

import "fmt"

type TestStruct struct {
	Field1 string
	Field2 int
}

func Function1() {}

func Function2() {}

func Function3() {}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReadCode(testFile, tt.functions...)
			if got != tt.want {
				t.Errorf("ReadCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddOrEditFunction(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")

	// Initial test content
	initialContent := `package main

import "fmt"

func ExistingFunction() {
	fmt.Println("Original content")
}
`
	if err := os.WriteFile(testFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test cases
	tests := []struct {
		name          string
		functionName  string
		functionBody  string
		want          string
		expectedError bool
	}{
		{
			name:         "Edit existing function",
			functionName: "ExistingFunction",
			functionBody: `func ExistingFunction() {
	fmt.Println("Modified content")
}`,
			want: `package main

import "fmt"

func ExistingFunction() {
	fmt.Println("Modified content")
}
`,
			expectedError: false,
		},
		{
			name:         "Add new function",
			functionName: "NewFunction",
			functionBody: `func NewFunction() {
	fmt.Println("New function content")
}`,
			want: `package main

import "fmt"

func ExistingFunction() {
	fmt.Println("Original content")
}

func NewFunction() {
	fmt.Println("New function content")
}
`,
			expectedError: false,
		},
		{
			name:          "Invalid function body",
			functionName:  "InvalidFunction",
			functionBody:  "invalid go code",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddOrEditFunction(testFile, tt.functionName, tt.functionBody)

			if tt.expectedError {
				if !strings.Contains(result, "Error") {
					t.Errorf("Expected error but got: %v", result)
				}
				return
			}

			if result != "Function successfully added/edited" {
				t.Errorf("Unexpected result: %v", result)
			}

			// Read the file to verify the changes
			content, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			if string(content) != tt.want {
				t.Errorf("File content = %v, want %v", string(content), tt.want)
			}
		})
	}
}
