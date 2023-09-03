package glox

import (
	"fmt"
	"strings"
)

type Expr any

func parenthesize(name string, exprs ...Expr) string {
	builder := &strings.Builder{}
	builder.WriteByte('(')
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteByte(' ')
		builder.WriteString(ExprToString(expr))
	}
	builder.WriteByte(')')
	return builder.String()
}

func ExprToString(e Expr) string {
	switch v := e.(type) {
	case UnaryExpr:
		return parenthesize(v.operator.Lexeme, v.right)
	case BinaryExpr:
		return parenthesize(v.operator.Lexeme, v.left, v.right)
	case Literal:
		if v.value == nil {
			return "nil"
		}
		return fmt.Sprint(v.value)
	case Grouping:
		return parenthesize("group", v.expr)
	}
	return fmt.Sprintf("unknown expr type: %v", e)
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

type Literal struct {
	value interface{}
}

type Grouping struct {
	expr Expr
}
