package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aselhid/indoscript/internal/ast"
	"github.com/aselhid/indoscript/internal/interpreter"
	"github.com/aselhid/indoscript/internal/lexer"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: indoscript [script].indos")
		os.Exit(1)
	}

	filename := os.Args[1]
	if err := runFile(filename); err != nil {
		reportError(err)
		os.Exit(1) // TODO: return exit code based on err
	}
}

func runFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	run(file)
	return nil
}

func run(reader io.Reader) {
	scanner := lexer.NewScanner(reader, os.Stderr)
	tokens := scanner.ScanTokens()
	parser := ast.NewParser(tokens)
	stmts, hasError := parser.Parse()
	if hasError {
		os.Exit(1)
	}
	interpreter := interpreter.NewInterpreter(os.Stdout, os.Stderr)
	interpreter.Interpret(stmts)
}
