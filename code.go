package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
)

func Lint(path string) string {
	path = Path(path)
	dir := filepath.Dir(path)

	command := exec.Command("go", "mod", "tidy")
	command.Dir = dir
	output, err := command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return fmt.Sprintf("Error formatting go file: %s", err)
	}

	command = exec.Command("go", "vet", dir)
	output, err = command.CombinedOutput()
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

// This function reads the code of the specified functions from the specified file.
// It also returns the rest of the function signatures (without the body) in the file
// and the structs, interfaces and types in the file.
// So, in summary, it returns the Go file but removing the body of any functions not passed in the functions parameter.
func ReadCode(path string, functions ...string) string {
	path = Path(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %s", err)
	}

	// Parse the Go file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return fmt.Sprintf("Error parsing file: %s", err)
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

	return result.String()
}

// This function adds a new function to the specified Go file, or edits an existing function.
func AddOrEditFunction(path string, functionName string, functionBody string) string {
	path = Path(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("Error reading file: %s", err)
	}

	// Parse the Go file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err != nil {
		return fmt.Sprintf("Error parsing file: %s", err)
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
		return fmt.Sprintf("Error parsing new function: %s", err)
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
		return "Error: Could not parse new function declaration"
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
		return fmt.Sprintf("Error writing file: %s", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Sprintf("Error saving file: %s", err)
	}

	return "Function successfully added/edited"
}
