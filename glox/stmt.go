package glox

import (
	"fmt"
	"strings"
)

type Stmt interface {
	Execute(env *Environment)
}

func StmtToString(s Stmt) string {
	switch v := s.(type) {
	case PrintStmt:
		return parenthesize("print", v.expr)
	case ExprStmt:
		return parenthesize("expr", v.expr)
	case VarDecl:
		return parenthesize("var "+v.name.Lexeme, v.initializer)
	case Block:
		stmtStrs := make([]string, 0, len(v.statements))
		for _, stmt := range v.statements {
			stmtStrs = append(stmtStrs, StmtToString(stmt))
		}
		return "{" + strings.Join(stmtStrs, "\n") + "}"
	case IfStmt:
		return parenthesize(
			parenthesize("if", v.condition) + "\n\t" +
				StmtToString(v.thenBranch) + "\n\t" +
				StmtToString(v.elseBranch),
		)
	case WhileStmt:
		return parenthesize(
			parenthesize("while", v.condition) + "\n\t" +
				StmtToString(v.body),
		)
	}
	return fmt.Sprintf("unknown stmt type: %v", s)
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
	name        Token
	initializer Expr
}

func (e VarDecl) Execute(env *Environment) {
	var v any = nil
	if e.initializer != nil {
		v = e.initializer.Evaluate(env)
	}
	if err := env.Declare(e.name.Lexeme, v); err != nil {
		panic(err)
	}
}

type Block struct {
	statements []Stmt
}

func (b Block) Execute(env *Environment) {
	newEnv := NewEnvironment(env)
	for _, stmt := range b.statements {
		stmt.Execute(newEnv)
	}
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func (s IfStmt) Execute(env *Environment) {
	result := s.condition.Evaluate(env)
	if isTruthy(result) {
		s.thenBranch.Execute(env)
	} else if s.elseBranch != nil {
		s.elseBranch.Execute(env)
	}
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (s WhileStmt) Execute(env *Environment) {
	for isTruthy(s.condition.Evaluate(env)) {
		s.body.Execute(env)
	}
}
