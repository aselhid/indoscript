package errors

import (
	"fmt"

	"github.com/aselhid/indoscript/internal/ast"
)

type RuntimeError struct {
	message string
	token   ast.Token
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("'%s' - %s", e.token.Lexeme, e.message)
}

func NewRuntimeError(token ast.Token, message string) error {
	return RuntimeError{token: token, message: message}
}
