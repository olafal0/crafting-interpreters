package glox

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenTypeEOF || p.current >= len(p.tokens)-1
}

// peek returns the current token without consuming it
func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// advance consumes the current token and returns it
func (p *Parser) advance() Token {
	t := p.tokens[p.current]
	if !p.isAtEnd() {
		p.current++
	}
	return t
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// check returns true if the next token is of the given type
func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

// match consumes a token as long as the type is one of the provided types
func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// consume consumes a token, as long as it is of the given type
func (p *Parser) consume(t TokenType) error {
	next := p.peek()
	if next.Type == t {
		p.advance()
		return nil
	}
	return fmt.Errorf("expected token type %v, got %v (line %d:%d)", t, next.Type, next.Pos.Line, next.Pos.Start)
}

func (p *Parser) Program() []Stmt {
	stmts := []Stmt{}
	for !p.isAtEnd() {
		stmts = append(stmts, p.Decl())
	}
	return stmts
}

func (p *Parser) Decl() Stmt {
	if p.match(TokenTypeVar) {
		return p.VarDecl()
	}
	return p.Statement()
}

func (p *Parser) VarDecl() Stmt {
	identifier := p.advance()
	var initializer Expr
	if p.match(TokenTypeEqual) {
		initializer = p.Expression()
	}
	if err := p.consume(TokenTypeSemicolon); err != nil {
		panic(err)
	}
	return VarDecl{name: identifier.Lexeme, initializer: initializer}
}

func (p *Parser) Statement() Stmt {
	if p.match(TokenTypePrint) {
		return p.PrintStmt()
	}
	return p.ExprStmt()
}

func (p *Parser) ExprStmt() Stmt {
	expr := p.Expression()
	err := p.consume(TokenTypeSemicolon)
	if err != nil {
		panic(err)
	}
	return ExprStmt{expr: expr}
}

func (p *Parser) PrintStmt() Stmt {
	expr := p.Expression()
	err := p.consume(TokenTypeSemicolon)
	if err != nil {
		panic(err)
	}
	return PrintStmt{expr: expr}
}

func (p *Parser) Expression() Expr {
	// Top-level short circuit: discard comments
	for p.check(TokenTypeComment) {
		p.advance()
	}
	return p.Equality()
}

func (p *Parser) Equality() Expr {
	expr := p.Comparison()
	for p.match(TokenTypeBangEqual, TokenTypeEqualEqual) {
		operator := p.previous()
		right := p.Comparison()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) Comparison() Expr {
	expr := p.Term()
	for p.match(TokenTypeGreater, TokenTypeEqual, TokenTypeLess, TokenTypeLessEqual) {
		operator := p.previous()
		right := p.Term()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) Term() Expr {
	expr := p.Factor()
	for p.match(TokenTypeMinus, TokenTypePlus) {
		operator := p.previous()
		right := p.Factor()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) Factor() Expr {
	expr := p.Unary()
	for p.match(TokenTypeSlash, TokenTypeStar) {
		operator := p.previous()
		right := p.Unary()
		expr = BinaryExpr{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) Unary() Expr {
	if p.match(TokenTypeBang, TokenTypeMinus) {
		operator := p.previous()
		right := p.Unary()
		return UnaryExpr{operator: operator, right: right}
	}
	return p.Primary()
}

func (p *Parser) Primary() Expr {
	switch {
	case p.match(TokenTypeFalse):
		return Literal{value: false}
	case p.match(TokenTypeTrue):
		return Literal{value: true}
	case p.match(TokenTypeNil):
		return Literal{value: nil}
	case p.match(TokenTypeString):
		return Literal{value: p.previous().Literal}
	case p.match(TokenTypeNumber):
		nStr, ok := p.previous().Literal.(string)
		if !ok {
			panic(fmt.Errorf("expected unparsed number, got %v (line %d:%d)", p.previous().Literal, p.previous().Pos.Line, p.previous().Pos.Start))
		}
		n, err := strconv.ParseFloat(nStr, 64)
		if err != nil {
			panic(fmt.Errorf("expected number, got %v (line %d:%d)", p.previous().Literal, p.previous().Pos.Line, p.previous().Pos.Start))
		}
		return Literal{value: n}
	case p.match(TokenTypeLeftParen):
		expr := p.Expression()
		err := p.consume(TokenTypeRightParen)
		if err != nil {
			// TODO: sync to next statement boundary (semicolon, keyword)
			panic(err)
		}
		return Grouping{expr: expr}
	case p.match(TokenTypeIdentifier):
		return Identifier{name: p.previous().Lexeme}
	}
	panic(fmt.Errorf("expected expression, got %v (line %d:%d)", p.peek().Type, p.peek().Pos.Line, p.peek().Pos.Start))
}

func (p *Parser) Execute(env *Environment) {
	statements := p.Program()
	for _, stmt := range statements {
		stmt.Execute(env)
	}
}
