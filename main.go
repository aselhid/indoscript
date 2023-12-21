package main

import (
	"fmt"
	"io"
	"os"

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
	if err := run(file); err != nil {
		return err
	}
	return nil
}

// run code coming from io.Reader
// used to run code from either file or stdin
func run(reader io.Reader) error {
	tokens, err := lexer.ScanTokens(reader)
	if err != nil {
		return err
	}
	return nil
}
