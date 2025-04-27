package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListDirectory(t *testing.T) {
	type args struct {
		path  string
		depth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty directory",
			args: args{
				path:  "testdata/empty",
				depth: 1,
			},
			want: "Empty directory",
		},
		{
			name: "single level directory",
			args: args{
				path:  "testdata/single",
				depth: 1,
			},
			want: "file1.txt\nfile2.txt",
		},
		{
			name: "nested directory with depth 1",
			args: args{
				path:  "testdata/nested",
				depth: 1,
			},
			want: "file1.txt\nsubdir",
		},
		{
			name: "nested directory with depth 2",
			args: args{
				path:  "testdata/nested",
				depth: 2,
			},
			want: "file1.txt\nsubdir\nsubdir/file2.txt",
		},
		{
			name: "directory with hidden files",
			args: args{
				path:  "testdata/hidden",
				depth: 1,
			},
			want: "visible.txt",
		},
		{
			name: "directory with gitignore",
			args: args{
				path:  "testdata/gitignore",
				depth: 1,
			},
			want: "included.txt",
		},
		{
			name: "negative depth",
			args: args{
				path:  "testdata/single",
				depth: -1,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListDirectory(tt.args.path, tt.args.depth); got != tt.want {
				t.Errorf("ListDirectory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGitignore(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
	}{
		{
			name:     "non-existent gitignore",
			path:     "testdata/nonexistent/.gitignore",
			patterns: []string{},
		},
		{
			name:     "empty gitignore",
			path:     "testdata/empty/.gitignore",
			patterns: []string{},
		},
		{
			name:     "gitignore with patterns",
			path:     "testdata/gitignore/.gitignore",
			patterns: []string{"ignored.txt"},
		},
		{
			name:     "gitignore with comments and empty lines",
			path:     "testdata/gitignore_complex/.gitignore",
			patterns: []string{"*.log", "build/", "temp.txt"},
		},
	}

	// Create test files
	err := os.WriteFile("testdata/empty/.gitignore", []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll("testdata/gitignore_complex", 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("testdata/gitignore_complex/.gitignore", []byte("# Ignore log files\n*.log\n\n# Ignore build directory\nbuild/\n\ntemp.txt"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGitignore(tt.path)
			if len(got.Patterns) != len(tt.patterns) {
				t.Errorf("NewGitignore() got %v patterns, want %v patterns", len(got.Patterns), len(tt.patterns))
			}
			for i, pattern := range tt.patterns {
				if i >= len(got.Patterns) || got.Patterns[i] != pattern {
					t.Errorf("NewGitignore() pattern %d = %v, want %v", i, got.Patterns[i], pattern)
				}
			}
		})
	}
}

func TestIsIgnored(t *testing.T) {
	gitignore := &Gitignore{
		Patterns: []string{"*.log", "build/", "temp.txt", "ignored*.txt"},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"should ignore log file", "app.log", true},
		{"should ignore nested log file", "logs/error.log", true},
		{"should ignore build directory", "build/output", true},
		{"should ignore specific file", "temp.txt", true},
		{"should ignore pattern match", "ignored-1.txt", true},
		{"should not ignore regular file", "main.go", false},
		{"should not ignore different extension", "log.txt", false},
		{"should not ignore partial match", "nottemp.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gitignore.IsIgnored(tt.path); got != tt.expected {
				t.Errorf("IsIgnored(%v) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	// Create test file
	testContent := "line1\nline2\nline3\nline4\nline5"
	err := os.WriteFile("testdata/read_test.txt", []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("testdata/read_test.txt")

	tests := []struct {
		name   string
		path   string
		offset int
		length int
		want   string
	}{
		{
			name:   "read entire file",
			path:   "testdata/read_test.txt",
			offset: 0,
			length: 0,
			want:   testContent,
		},
		{
			name:   "read with offset",
			path:   "testdata/read_test.txt",
			offset: 2,
			length: 2,
			want:   "line3\nline4",
		},
		{
			name:   "read non-existent file",
			path:   "testdata/nonexistent.txt",
			offset: 0,
			length: 1,
			want:   "Error reading file: open " + filepath.Join(workingDirectory, "testdata/nonexistent.txt") + ": no such file or directory",
		},
		{
			name:   "read with invalid offset",
			path:   "testdata/read_test.txt",
			offset: 10,
			length: 1,
			want:   "File has 5 lines, cannot read line 10",
		},
		{
			name:   "read with large length",
			path:   "testdata/read_test.txt",
			offset: 0,
			length: 1001,
			want:   "Cannot read more than 1000 lines",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadFile(tt.path, tt.offset, tt.length); got != tt.want {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		content string
		wantErr bool
	}{
		{
			name:    "write to new file",
			path:    "testdata/write_test.txt",
			content: "test content",
			wantErr: false,
		},
		{
			name:    "write to existing file",
			path:    "testdata/write_test.txt",
			content: "new content",
			wantErr: false,
		},
		{
			name:    "write to invalid path",
			path:    "/invalid/path/file.txt",
			content: "test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WriteFile(tt.path, tt.content)
			if tt.wantErr {
				if !strings.Contains(got, "Error writing to file:") {
					t.Errorf("WriteFile() error = %v, wantErr %v", got, tt.wantErr)
				}
			} else {
				// Verify file was written correctly
				content, err := os.ReadFile(Path(tt.path))
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
				}
				if string(content) != tt.content {
					t.Errorf("WriteFile() wrote %v, want %v", string(content), tt.content)
				}
			}
		})
	}

	// Cleanup
	os.Remove("testdata/write_test.txt")
}

func TestMkDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create new directory",
			path:    "testdata/new_dir",
			wantErr: false,
		},
		{
			name:    "create nested directories",
			path:    "testdata/parent/child/grandchild",
			wantErr: false,
		},
		{
			name:    "create existing directory",
			path:    "testdata/new_dir",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MkDir(tt.path)
			if tt.wantErr {
				if !strings.Contains(got, "Error creating directory:") {
					t.Errorf("MkDir() error = %v, wantErr %v", got, tt.wantErr)
				}
			} else {
				// Verify directory was created
				if _, err := os.Stat(Path(tt.path)); os.IsNotExist(err) {
					t.Errorf("MkDir() failed to create directory %v", tt.path)
				}
			}
		})
	}

	// Cleanup
	os.RemoveAll("testdata/new_dir")
	os.RemoveAll("testdata/parent")
}

func TestPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "empty path",
			path: "",
			want: workingDirectory,
		},
		{
			name: "current directory",
			path: ".",
			want: workingDirectory,
		},
		{
			name: "relative path",
			path: "testdata/file.txt",
			want: filepath.Join(workingDirectory, "testdata/file.txt"),
		},
		{
			name: "absolute path",
			path: "/absolute/path/file.txt",
			want: "/absolute/path/file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Path(tt.path); got != tt.want {
				t.Errorf("Path() = %v, want %v", got, tt.want)
			}
		})
	}
}
