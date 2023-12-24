package interpreter

import (
	"github.com/aselhid/indoscript/internal/ast"
	"github.com/aselhid/indoscript/internal/environment"
)

type Callable interface {
	Call(*Interpreter, []any) any
}

type FunctionCallable struct {
	Declaration ast.FuncStmt
}

type ReturnValue struct {
	Value any
}

func (f *FunctionCallable) Call(interpreter *Interpreter, arguments []any) (returnValue any) {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(ReturnValue); ok {
				returnValue = v.Value
				return
			}
			panic(err)
		}
	}()

	env := environment.NewEnvironment(interpreter.globalEnv)
	for i, declaration := range f.Declaration.Parameters {
		env.Define(declaration, arguments[i])
	}
	interpreter.executeBlock(f.Declaration.Body, env)
	return nil
}

func NewFunctionCallable(declaration ast.FuncStmt) FunctionCallable {
	return FunctionCallable{
		Declaration: declaration,
	}
}
