package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDirectory(t *testing.T) {
	// Create a temporary test directory
	testDir, err := os.MkdirTemp("", "listdir_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	// Create test directory structure
	dirs := []string{
		"dir1",
		"dir1/subdir1",
		"dir2",
	}
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(testDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Create test files
	files := []string{
		"file1.txt",
		"dir1/file2.txt",
		"dir1/subdir1/file3.txt",
		"dir2/file4.txt",
	}
	for _, file := range files {
		err := os.WriteFile(filepath.Join(testDir, file), []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test cases
	tests := []struct {
		name     string
		path     string
		depth    int
		expected string
	}{
		{
			name:     "depth 0 should return empty string",
			path:     testDir,
			depth:    0,
			expected: "",
		},
		{
			name:     "depth 1 should list only top level",
			path:     testDir,
			depth:    1,
			expected: "dir1\ndir2\nfile1.txt",
		},
		{
			name:     "depth 2 should list two levels",
			path:     testDir,
			depth:    2,
			expected: "dir1\ndir1/file2.txt\ndir1/subdir1\ndir2\ndir2/file4.txt\nfile1.txt",
		},
		{
			name:     "depth 3 should list all levels",
			path:     testDir,
			depth:    3,
			expected: "dir1\ndir1/file2.txt\ndir1/subdir1\ndir1/subdir1/file3.txt\ndir2\ndir2/file4.txt\nfile1.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ListDirectory(tt.path, tt.depth)
			if result != tt.expected {
				t.Errorf("ListDirectory(%q, %d) = %q, want %q", tt.path, tt.depth, result, tt.expected)
			}
		})
	}
}
