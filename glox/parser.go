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
	return fmt.Errorf("expected %s token, got %s", t, next)
}

func (p *Parser) Program() []Stmt {
	stmts := []Stmt{}
	for !p.isAtEnd() {
		// Top-level short circuit: discard comments
		for p.check(TokenTypeComment) {
			p.advance()
		}
		if p.isAtEnd() {
			break
		}
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
	return VarDecl{name: identifier, initializer: initializer}
}

func (p *Parser) Statement() Stmt {
	switch {
	case p.match(TokenTypeIf):
		return p.IfStmt()
	case p.match(TokenTypePrint):
		return p.PrintStmt()
	case p.match(TokenTypeLeftBrace):
		return p.Block()
	case p.match(TokenTypeWhile):
		return p.WhileStmt()
	case p.match(TokenTypeFor):
		return p.ForStmt()
	}
	return p.ExprStmt()
}

func (p *Parser) IfStmt() Stmt {
	// Don't consume left paren here, to support if statements without parens
	condition := p.Expression()
	// Don't consume right paren

	thenBranch := p.Statement()
	var elseBranch Stmt
	if p.match(TokenTypeElse) {
		elseBranch = p.Statement()
	}
	return IfStmt{condition: condition, thenBranch: thenBranch, elseBranch: elseBranch}
}

func (p *Parser) Block() Stmt {
	statements := []Stmt{}
	for !p.check(TokenTypeRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.Decl())
	}
	if err := p.consume(TokenTypeRightBrace); err != nil {
		panic(err)
	}
	return Block{statements: statements}
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

func (p *Parser) WhileStmt() Stmt {
	// consume left paren in case one's provided
	p.match(TokenTypeLeftParen)
	condition := p.Expression()
	// consume left paren in case one's provided
	// TODO: complain if there was no left paren
	p.match(TokenTypeRightParen)
	body := p.Statement()
	return WhileStmt{condition: condition, body: body}
}

func (p *Parser) ForStmt() Stmt {
	// consume left paren in case one's provided
	p.match(TokenTypeLeftParen)
	var initializer Stmt
	if p.match(TokenTypeSemicolon) {
		initializer = nil
	} else if p.match(TokenTypeVar) {
		initializer = p.VarDecl()
	} else {
		initializer = p.ExprStmt()
	}

	var condition Expr
	if !p.check(TokenTypeSemicolon) {
		condition = p.Expression()
	}

	var increment Expr
	if p.match(TokenTypeSemicolon) {
		increment = p.Expression()
	}

	// consume left paren in case one's provided
	// TODO: complain if there was no left paren
	p.match(TokenTypeRightParen)

	body := p.Statement()
	if increment != nil {
		body = Block{
			statements: []Stmt{
				body,
				ExprStmt{
					expr: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = Literal{
			value: true,
		}
	}
	body = WhileStmt{
		condition: condition,
		body:      body,
	}

	if initializer != nil {
		body = Block{
			statements: []Stmt{
				initializer,
				body,
			},
		}
	}
	return body
}

func (p *Parser) Expression() Expr {
	return p.Assignment()
}

func (p *Parser) Assignment() Expr {
	expr := p.LogicOr()
	if p.match(TokenTypeEqual) {
		val := p.Assignment()
		if exprVar, ok := expr.(Identifier); ok {
			name := exprVar.name
			return Assign{name: name, val: val}
		}
		panic(fmt.Errorf("invalid assign target %s (%s)", ExprToString(expr), expr.Pos()))
	}
	return expr
}

func (p *Parser) LogicOr() Expr {
	expr := p.LogicAnd()
	for p.match(TokenTypeOr) {
		operator := p.previous()
		right := p.LogicAnd()
		expr = Logical{left: expr, operator: operator, right: right}
	}
	return expr
}

func (p *Parser) LogicAnd() Expr {
	expr := p.Equality()
	for p.match(TokenTypeAnd) {
		operator := p.previous()
		right := p.Equality()
		expr = Logical{left: expr, operator: operator, right: right}
	}
	return expr
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
	for p.match(TokenTypeGreater, TokenTypeEqualEqual, TokenTypeLess, TokenTypeLessEqual) {
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
		return Literal{token: p.previous(), value: false}
	case p.match(TokenTypeTrue):
		return Literal{token: p.previous(), value: true}
	case p.match(TokenTypeNil):
		return Literal{token: p.previous(), value: nil}
	case p.match(TokenTypeString):
		return Literal{token: p.previous(), value: p.previous().Literal}
	case p.match(TokenTypeNumber):
		nStr, ok := p.previous().Literal.(string)
		if !ok {
			panic(fmt.Errorf("expected unparsed number, got %s", p.previous()))
		}
		n, err := strconv.ParseFloat(nStr, 64)
		if err != nil {
			panic(fmt.Errorf("expected number, got %s", p.previous()))
		}
		return Literal{token: p.previous(), value: n}
	case p.match(TokenTypeLeftParen):
		leftParen := p.previous()
		expr := p.Expression()
		err := p.consume(TokenTypeRightParen)
		if err != nil {
			// TODO: sync to next statement boundary (semicolon, keyword)
			panic(err)
		}
		return Grouping{left: leftParen, expr: expr, right: p.previous()}
	case p.match(TokenTypeIdentifier):
		return Identifier{name: p.previous()}
	}
	panic(fmt.Errorf("expected expression, got %s", p.peek()))
}

func (p *Parser) Execute(env *Environment) {
	statements := p.Program()
	for _, stmt := range statements {
		stmt.Execute(env)
	}
}

func (p *Parser) PrintAST() {
	statements := p.Program()
	for _, stmt := range statements {
		fmt.Println(StmtToString(stmt))
	}
}
