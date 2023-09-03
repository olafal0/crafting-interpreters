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

	TokenNames = map[TokenType]string{
		TokenTypeNone:         "none",
		TokenTypeLeftParen:    "leftparen",
		TokenTypeRightParen:   "rightparen",
		TokenTypeLeftBrace:    "leftbrace",
		TokenTypeRightBrace:   "rightbrace",
		TokenTypeComma:        "comma",
		TokenTypeDot:          "dot",
		TokenTypeMinus:        "minus",
		TokenTypePlus:         "plus",
		TokenTypeSemicolon:    "semicolon",
		TokenTypeSlash:        "slash",
		TokenTypeStar:         "star",
		TokenTypeBang:         "bang",
		TokenTypeBangEqual:    "bangequal",
		TokenTypeEqual:        "equal",
		TokenTypeEqualEqual:   "equalequal",
		TokenTypeGreater:      "greater",
		TokenTypeGreaterEqual: "greaterequal",
		TokenTypeLess:         "less",
		TokenTypeLessEqual:    "lessequal",
		TokenTypeIdentifier:   "identifier",
		TokenTypeString:       "string",
		TokenTypeNumber:       "number",
		TokenTypeComment:      "comment",
		TokenTypeAnd:          "and",
		TokenTypeClass:        "class",
		TokenTypeElse:         "else",
		TokenTypeFalse:        "false",
		TokenTypeFun:          "fun",
		TokenTypeFor:          "for",
		TokenTypeIf:           "if",
		TokenTypeNil:          "nil",
		TokenTypeOr:           "or",
		TokenTypePrint:        "print",
		TokenTypeReturn:       "return",
		TokenTypeSuper:        "super",
		TokenTypeThis:         "this",
		TokenTypeTrue:         "true",
		TokenTypeVar:          "var",
		TokenTypeWhile:        "while",
		TokenTypeEOF:          "eof",
	}
)

func (t TokenType) String() string {
	if name, ok := TokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", t)
}

type Pos struct {
	Line  int
	Start int
	End   int
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Pos     Pos
}

func (t Token) String() string {
	return fmt.Sprintf("Line %d [%d:%d] Type %s '%s' %v", t.Pos.Line, t.Pos.Start, t.Pos.End, t.Type, t.Lexeme, t.Literal)
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
