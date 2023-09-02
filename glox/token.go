package glox

import "fmt"

type TokenType int16

const (
	TokenTypeNone TokenType = iota

	// Single-character tokens
	TokenTypeLeftParen
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemicolon
	TokenTypeSlash
	TokenTypeStar

	// One or two character tokens
	TokenTypeBang
	TokenTypeBangEqual
	TokenTypeEqual
	TokenTypeEqualEqual
	TokenTypeGreater
	TokenTypeGreaterEqual
	TokenTypeLess
	TokenTypeLessEqual

	// Literals
	TokenTypeIdentifier
	TokenTypeString
	TokenTypeNumber
	TokenTypeComment

	// Keywords
	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFun
	TokenTypeFor
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint // TODO: this should be a stdlib function call
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeTrue
	TokenTypeVar
	TokenTypeWhile

	TokenTypeEOF
)

var (
	ReservedKeywords = map[string]TokenType{
		"and":    TokenTypeAnd,
		"class":  TokenTypeClass,
		"else":   TokenTypeElse,
		"false":  TokenTypeFalse,
		"fun":    TokenTypeFun,
		"for":    TokenTypeFor,
		"if":     TokenTypeIf,
		"nil":    TokenTypeNil,
		"or":     TokenTypeOr,
		"print":  TokenTypePrint,
		"return": TokenTypeReturn,
		"super":  TokenTypeSuper,
		"this":   TokenTypeThis,
		"true":   TokenTypeTrue,
		"var":    TokenTypeVar,
		"while":  TokenTypeWhile,
	}
)

type Pos struct {
	Line  int32
	Start int32
	End   int32
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Pos     Pos
}

func (t Token) String() string {
	return fmt.Sprintf("Line %d [%d:%d] Type %d '%s' %v", t.Pos.Line, t.Pos.Start, t.Pos.End, t.Type, t.Lexeme, t.Literal)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
