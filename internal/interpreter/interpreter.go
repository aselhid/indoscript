package interpreter

import (
	"fmt"
	"io"

	"github.com/aselhid/indoscript/internal/ast"
	"github.com/aselhid/indoscript/internal/environment"
	"github.com/aselhid/indoscript/internal/errors"
	"github.com/google/go-cmp/cmp"
)

var globalEnvironment = environment.NewEnvironment(nil)

type Interpreter struct {
	stdErr io.Writer
	stdOut io.Writer
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) (hasRuntimeError bool) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(errors.RuntimeError); ok {
				i.stdErr.Write([]byte(e.Error() + "\n"))
			} else {
				fmt.Printf("Error: %s\n", err)
			}
			hasRuntimeError = true
		}
	}()

	for _, stmt := range stmts {
		i.execute(stmt)
	}
	return false
}

func (i *Interpreter) VisitVarStmt(stmt ast.VarStmt) {
	value := i.evaluate(stmt.Expression)
	globalEnvironment.Define(stmt.Identifier, value)
}

func (i *Interpreter) VisitPrintStmt(stmt ast.PrintStmt) {
	expr := i.evaluate(stmt.Expression)
	i.stdOut.Write([]byte(i.stringify(expr) + "\n"))
}

func (i *Interpreter) VisitExprStmt(stmt ast.ExprStmt) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitBlockStmt(stmt ast.BlockStmt) {
	env := environment.NewEnvironment(globalEnvironment)
	i.executeBlock(stmt.Statements, env)
}

func (i *Interpreter) VisitBinaryExpr(expr ast.BinaryExpr) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case ast.TokenMinus:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case ast.TokenSlash:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case ast.TokenStar:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	case ast.TokenPlus:
		leftAsNumber, leftIsNumber := left.(float64)
		rightAsNumber, rightIsNumber := right.(float64)
		if leftIsNumber && rightIsNumber {
			return leftAsNumber + rightAsNumber
		}

		leftAsString, leftIsString := left.(string)
		rightAsString, rightIsString := right.(string)
		if leftIsString && rightIsString {
			return leftAsString + rightAsString
		}

		i.error(expr.Operator, "operands must be either numbers or strings")
	case ast.TokenGreater:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case ast.TokenGreaterEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case ast.TokenLess:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case ast.TokenLessEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case ast.TokenAnd:
		return i.isTruthy(left) && i.isTruthy(right)
	case ast.TokenOr:
		return i.isTruthy(left) || i.isTruthy(right)
	case ast.TokenEqualEqual:
		return cmp.Equal(left, right)
	case ast.TokenBangEqual:
		return !cmp.Equal(left, right)
	}
	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr ast.UnaryExpr) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case ast.TokenMinus:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case ast.TokenBang:
		return !i.isTruthy(right)
	}
	return nil
}

func (i *Interpreter) VisitPrimaryExpr(expr ast.PrimaryExpr) any {
	return expr.Literal
}

func (i *Interpreter) VisitGroupExpr(expr ast.GroupExpr) any {
	return i.evaluate(expr)
}

func (i *Interpreter) VisitVarExpr(expr ast.VarExpr) any {
	return globalEnvironment.Get(expr.Identifier)
}

func (i *Interpreter) isTruthy(value any) bool {
	switch v := value.(type) {
	case float64:
		return v != 0.0
	case string:
		return v != ""
	case bool:
		return v
	default:
		return false
	}
}

func (i *Interpreter) checkNumberOperand(token ast.Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}
	i.error(token, "operand must be a number")
}

func (i *Interpreter) checkNumberOperands(token ast.Token, left, right any) {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return
		}
	}
	i.error(token, "operands must be numbers")
}

func (i *Interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) executeBlock(stmts []ast.Stmt, environment *environment.Environment) {
	previousEnv := globalEnvironment
	globalEnvironment = environment
	for _, stmt := range stmts {
		i.execute(stmt)
	}
	globalEnvironment = previousEnv
}

func (i *Interpreter) error(token ast.Token, message string) {
	err := errors.NewRuntimeError(token, message)
	i.stdErr.Write([]byte(fmt.Sprintf("[line %d] Runtime error: %s\n", token.LineNumber, err.Error())))
	panic(err)
}

func (i *Interpreter) stringify(value any) string {
	switch v := value.(type) {
	case float64:
		return fmt.Sprintf("%f", v)
	case string:
		return v
	case bool:
		if v {
			return "benar"
		}
		return "salah"
	}
	return "unknown value"
}

func NewInterpreter(stdOut, stdErr io.Writer) *Interpreter {
	return &Interpreter{
		stdOut: stdOut,
		stdErr: stdErr,
	}
}
