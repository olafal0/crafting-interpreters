package glox

import (
	"fmt"
	"time"
)

type Caller interface {
	Call(env *Environment, args []any) any
	Arity() int
	String() string
}

type DefinedFunc struct {
	decl FuncDecl
}

var _ Caller = DefinedFunc{}

func (f DefinedFunc) Arity() int { return len(f.decl.params) }

func (f DefinedFunc) Call(env *Environment, args []any) any {
	funcEnv := NewEnvironment(env)
	for i := range f.decl.params {
		funcEnv.Declare(f.decl.params[i].Lexeme, args[i])
	}

	for _, stmt := range f.decl.body {
		stmt.Execute(funcEnv)
	}
	return nil
}

func (f DefinedFunc) String() string {
	return fmt.Sprintf("<fn %s>", f.decl.name.Lexeme)
}

type ClockFunc struct{}

var _ Caller = ClockFunc{}

func (f ClockFunc) Arity() int { return 0 }

func (f ClockFunc) Call(env *Environment, args []any) any {
	return float64(time.Now().UnixMicro()) / 1_000_000
}
func (f ClockFunc) String() string {
	return "<builtin fn clock>"
}
