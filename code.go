package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func Lint(path string) string {
	path = Path(path)

	extension := filepath.Ext(path)
	dir := filepath.Dir(path)

	switch extension {
	case ".go":
		command := exec.Command("go", "vet", dir)
		output, err := command.CombinedOutput()
		if err != nil && len(output) == 0 {
			return fmt.Sprintf("Error formatting go file: %s", err)
		}
		return string(output)
	}

	return "Linting not supported for this file type"
}
