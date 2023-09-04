package glox

import (
	"fmt"
	"strings"
)

type Expr interface {
	Evaluate(env *Environment) any
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

func (e UnaryExpr) Evaluate(env *Environment) any {
	right := e.right.Evaluate(env)
	switch e.operator.Type {
	case TokenTypeMinus:
		// TODO: check for runtime errors (can only negate numbers)
		return -right.(float64)
	case TokenTypeBang:
		return !isTruthy(right)
	}
	return nil
}

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func (e BinaryExpr) Evaluate(env *Environment) any {
	left := e.left.Evaluate(env)
	right := e.right.Evaluate(env)

	// TODO: check for runtime errors
	switch e.operator.Type {
	case TokenTypeMinus:
		return left.(float64) - right.(float64)
	case TokenTypeSlash:
		return left.(float64) / right.(float64)
	case TokenTypeStar:
		return left.(float64) * right.(float64)
	case TokenTypePlus:
		// Special case: we can add numbers or concatenate strings
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
		}
	case TokenTypeGreater:
		return left.(float64) > right.(float64)
	case TokenTypeGreaterEqual:
		return left.(float64) >= right.(float64)
	case TokenTypeLess:
		return left.(float64) < right.(float64)
	case TokenTypeLessEqual:
		return left.(float64) <= right.(float64)
	case TokenTypeBangEqual:
		return !isEqual(left, right)
	case TokenTypeEqualEqual:
		return isEqual(left, right)
	}
	return nil
}

type Literal struct {
	value interface{}
}

func (e Literal) Evaluate(env *Environment) any {
	return e.value
}

type Grouping struct {
	expr Expr
}

func (e Grouping) Evaluate(env *Environment) any {
	return e.expr.Evaluate(env)
}

type Identifier struct {
	name string
}

func (e Identifier) Evaluate(env *Environment) any {
	v, ok := env.Get(e.name)
	if !ok {
		panic(fmt.Errorf("unknown identifier %s", e.name))
	}
	return v
}

func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

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
