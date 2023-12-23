package environment

import (
	"fmt"

	"github.com/aselhid/indoscript/internal/ast"
	"github.com/aselhid/indoscript/internal/errors"
)

type Environment struct {
	values   map[string]any
	encloser *Environment
}

func (e *Environment) Define(identifier ast.Token, value any) {
	e.values[identifier.Lexeme] = value
}

func (e *Environment) Assign(identifier ast.Token, value any) {
	if _, ok := e.values[identifier.Lexeme]; !ok {
		if e.encloser == nil {
			e.error(identifier, fmt.Sprintf("Undefined variable %s\n", identifier.Lexeme))
		}
		e.encloser.Assign(identifier, value)
	}

	e.values[identifier.Lexeme] = value
}

func (e *Environment) Get(identifier ast.Token) any {
	value, ok := e.values[identifier.Lexeme]
	if !ok {
		if e.encloser == nil {
			e.error(identifier, fmt.Sprintf("Undefined variable %s\n", identifier.Lexeme))
		}
		return e.encloser.Get(identifier)
	}
	return value
}

func (e *Environment) error(token ast.Token, message string) {
	err := errors.NewRuntimeError(token, message)
	panic(err)
}

func NewEnvironment(encloser *Environment) *Environment {
	return &Environment{
		values:   make(map[string]any),
		encloser: encloser,
	}
}
