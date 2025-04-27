package main

import (
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
