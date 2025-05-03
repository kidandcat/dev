package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"

	"golang.org/x/tools/imports"
)

func Lint(path string) map[string]any {
	path = Path(path)
	dir := filepath.Dir(path)

	command := exec.Command("go", "mod", "tidy")
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return map[string]any{
			"error": fmt.Sprintf("Error formatting go file: %s", err),
		}
	}

	command = exec.Command("go", "vet", dir)
	output, err = command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return map[string]any{
			"error": fmt.Sprintf("Error formatting go file: %s", err),
		}
	}
	if len(output) > 0 {
		return map[string]any{
			"error": string(output),
		}
	}

	command = exec.Command("go", "fmt", dir)
	output, err = command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return map[string]any{
			"error": fmt.Sprintf("Error formatting go file: %s", err),
		}
	}
	if len(output) > 0 {
		return map[string]any{
			"error": string(output),
		}
	}

	autoImport(path)

	return map[string]any{
		"results": "No errors found",
	}
}

// This function reads the code of the specified functions from the specified file.
// It also returns the rest of the function signatures (without the body) in the file
// and the structs, interfaces and types in the file.
// So, in summary, it returns the Go file but removing the body of any functions not passed in the functions parameter.
func ReadCode(path string, functions ...string) map[string]any {
	path = Path(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error reading file: %s", err),
		}
	}

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Get the package name from the go.mod file
		mod, err := os.ReadFile(filepath.Join(filepath.Dir(path), "go.mod"))
		if err != nil {
			return map[string]any{
				"error": fmt.Sprintf("Error reading go.mod file: %s", err),
			}
		}
		modString := string(mod)
		packageName := strings.TrimSpace(strings.Split(modString, "\n")[0])
		// Create from template
		content = []byte(fmt.Sprintf("package %s\n\n", packageName))
		if err := os.WriteFile(path, content, 0644); err != nil {
			return map[string]any{
				"error": fmt.Sprintf("Error creating file: %s", err),
			}
		}
	}

	// Parse the Go file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error parsing file: %s", err),
		}
	}

	var result strings.Builder

	// Print package declaration
	result.WriteString(fmt.Sprintf("package %s\n\n", f.Name.Name))

	// Print imports
	if len(f.Imports) > 0 {
		result.WriteString("import (\n")
		for _, imp := range f.Imports {
			if imp.Name != nil {
				result.WriteString(fmt.Sprintf("\t%s %s\n", imp.Name.Name, imp.Path.Value))
			} else {
				result.WriteString(fmt.Sprintf("\t%s\n", imp.Path.Value))
			}
		}
		result.WriteString(")\n\n")
	}

	// Print declarations
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// Check if this function is in the list of functions to keep
			shouldKeep := false
			for _, fn := range functions {
				if d.Name.Name == fn {
					shouldKeep = true
					break
				}
			}

			if shouldKeep {
				// Print the complete function
				printer.Fprint(&result, fset, d)
				result.WriteString("\n\n")
			} else {
				// Print only the function signature
				result.WriteString("func ")
				if d.Recv != nil {
					printer.Fprint(&result, fset, d.Recv)
					result.WriteString(" ")
				}
				result.WriteString(d.Name.Name)
				printer.Fprint(&result, fset, d.Type)
				result.WriteString(" {}\n\n")
			}
		default:
			// Print all other declarations (types, vars, etc.)
			printer.Fprint(&result, fset, d)
			result.WriteString("\n\n")
		}
	}

	return map[string]any{
		"results": result.String(),
	}
}

// This function adds a new function to the specified Go file, or edits an existing function.
func AddOrEditFunction(path string, functionName string, functionBody string) map[string]any {
	path = Path(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error reading file: %s", err),
		}
	}

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File do not exists, create it
		os.Create(path)
		content, err = os.ReadFile(path)
		if err != nil {
			return map[string]any{
				"error": fmt.Sprintf("Error reading file: %s", err),
			}
		}
	}
	// Parse the Go file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error parsing file: %s", err),
		}
	}

	// Check if function already exists
	var existingFunc *ast.FuncDecl
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn.Name.Name == functionName {
			existingFunc = fn
			break
		}
	}

	// Parse the new function body
	newFunc, err := parser.ParseFile(fset, "", "package p\n"+functionBody, parser.ParseComments)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error parsing new function: %s", err),
		}
	}

	// Get the new function declaration
	var newFuncDecl *ast.FuncDecl
	for _, decl := range newFunc.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			newFuncDecl = fn
			break
		}
	}

	if newFuncDecl == nil {
		return map[string]any{
			"error": "Error: Could not parse new function declaration",
		}
	}

	// If function exists, replace it; otherwise, add it
	if existingFunc != nil {
		// Replace the existing function
		existingFunc.Body = newFuncDecl.Body
	} else {
		// Add the new function
		f.Decls = append(f.Decls, newFuncDecl)
	}

	// Write the modified file
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, f); err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error writing file: %s", err),
		}
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return map[string]any{
			"error": fmt.Sprintf("Error saving file: %s", err),
		}
	}

	// Lint the file
	lintResult := Lint(path)

	return map[string]any{
		"results": "Function successfully added/edited",
		"lint":    lintResult,
	}
}

func autoImport(path string) {
	filename := Path(path)
	src, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	processedSrc, err := imports.Process(filename, src, nil)
	if err != nil {
		log.Fatalf("Error processing file with goimports: %v", err)
	}

	// If the processed content is different from the original, write it back
	if !bytes.Equal(src, processedSrc) {
		err = os.WriteFile(filename, processedSrc, 0644)
		if err != nil {
			log.Fatalf("Error writing processed file: %v", err)
		}
		fmt.Printf("Processed and updated %s\n", filename)
	} else {
		fmt.Printf("%s is already correctly formatted and imported.\n", filename)
	}
}
