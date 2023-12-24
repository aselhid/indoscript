package interpreter

type Callable interface {
	Call(*Interpreter, []any) any
}

type Function struct {
}

// func (f *Function) Call(interpreter *Interpreter, arguments []any) any {
// 	env := environment.NewEnvironment(int)
// }
