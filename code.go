package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Lint(path string) string {
	path = Path(path)
	dir := filepath.Dir(path)

	command := exec.Command("go", "vet", dir)
	output, err := command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return fmt.Sprintf("Error formatting go file: %s", err)
	}
	if len(output) > 0 {
		return string(output)
	}

	command = exec.Command("go", "fmt", dir)
	output, err = command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return fmt.Sprintf("Error formatting go file: %s", err)
	}
	if len(output) > 0 {
		return string(output)
	}

	return "No errors found"
}
