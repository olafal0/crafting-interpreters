package glox

import "fmt"

type Stmt interface {
	Execute(env *Environment)
}

type PrintStmt struct {
	expr Expr
}

func (p PrintStmt) Execute(env *Environment) {
	fmt.Println(p.expr.Evaluate(env))
}

type ExprStmt struct {
	expr Expr
}

func (e ExprStmt) Execute(env *Environment) {
	e.expr.Evaluate(env)
}

type VarDecl struct {
	name        string
	initializer Expr
}

func (e VarDecl) Execute(env *Environment) {
	var v any = nil
	if e.initializer != nil {
		v = e.initializer.Evaluate(env)
	}
	if err := env.Declare(e.name, v); err != nil {
		panic(err)
	}
}
