package glox

import "testing"

func TestScanning(t *testing.T) {
	source := []byte(`// This is a comment
var x fun longIdentifier
123 / 12.3 + 0.123
"hello"
! = != > >=`)
	s := NewScanner(source)
	tokens, err := s.ScanTokens()
	if err != nil {
		t.Error(err)
	}
	expectedTypes := []TokenType{
		TokenTypeComment,
		TokenTypeVar,
		TokenTypeIdentifier,
		TokenTypeFun,
		TokenTypeIdentifier,
		TokenTypeNumber,
		TokenTypeSlash,
		TokenTypeNumber,
		TokenTypePlus,
		TokenTypeNumber,
		TokenTypeString,
		TokenTypeBang,
		TokenTypeEqual,
		TokenTypeBangEqual,
		TokenTypeGreater,
		TokenTypeGreaterEqual,
		TokenTypeEOF,
	}
	expectedLexemes := []string{
		"// This is a comment",
		"var",
		"x",
		"fun",
		"longIdentifier",
		"123",
		"/",
		"12.3",
		"+",
		"0.123",
		"\"hello\"",
		"!",
		"=",
		"!=",
		">",
		">=",
		"",
	}
	for i, token := range tokens {
		if token.Type != expectedTypes[i] {
			t.Errorf("Expected token type %v, got %v", expectedTypes[i], token.Type)
		}
		if token.Lexeme != expectedLexemes[i] {
			t.Errorf("Expected token lexeme %v, got %v", expectedLexemes[i], token.Lexeme)
		}
	}
}
