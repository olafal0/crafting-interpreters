package glox

import (
	"errors"
	"fmt"
)

type Scanner struct {
	current int
	line    int
	// visitedLinesLen is the number of characters in lines already visited. This
	// can be subtracted from current to get the position within a line.
	visitedLinesLen int
	source          []byte
	tokens          []Token
}

func NewScanner(source []byte) *Scanner {
	s := &Scanner{
		source: source,
		tokens: []Token{},
	}
	return s
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	errs := []error{}
	for int(s.current) < len(s.source) {
		err := s.scanToken()
		if err != nil {
			errs = append(errs, err)
		}
	}
	s.tokens = append(s.tokens, Token{
		Type:    TokenTypeEOF,
		Lexeme:  "",
		Literal: nil,
		Pos: Pos{
			Line:  s.line + 1,
			Start: 0,
			End:   0,
		},
	})
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	start := s.current
	c := s.source[s.current]
	s.current++
	var peek byte
	if s.current < len(s.source) {
		peek = s.source[s.current]
	}
	switch c {
	// Comments, newlines, and whitespace
	case '\n':
		// For newlines, ignore but increment the line counter
		s.line++
		s.visitedLinesLen = s.current
	case ' ', '\r', '\t':
	// Ignore whitespace
	case '/':
		if peek == '/' {
			// Consume characters until end of line
			for s.source[s.current] != '\n' && s.current < len(s.source) {
				s.current++
			}
			// Strip off double slash and leading space
			comment := string(s.source[start+2 : s.current])
			if len(comment) > 1 && comment[0] == ' ' {
				comment = comment[1:]
			}
			s.addLiteralToken(start, TokenTypeComment, comment)
			s.line++
			s.visitedLinesLen = s.current
		} else {
			s.addToken(start, TokenTypeSlash)
		}

	// Single-character tokens
	case '(':
		s.addToken(start, TokenTypeLeftParen)
	case ')':
		s.addToken(start, TokenTypeRightParen)
	case '{':
		s.addToken(start, TokenTypeLeftBrace)
	case '}':
		s.addToken(start, TokenTypeRightBrace)
	case ',':
		s.addToken(start, TokenTypeComma)
	case '.':
		s.addToken(start, TokenTypeDot)
	case '-':
		s.addToken(start, TokenTypeMinus)
	case '+':
		s.addToken(start, TokenTypePlus)
	case ';':
		s.addToken(start, TokenTypeSemicolon)
	case '*':
		s.addToken(start, TokenTypeStar)

		// 1-2 character tokens
	case '!':
		if peek == '=' {
			s.current++
			s.addToken(start, TokenTypeBangEqual)
		} else {
			s.addToken(start, TokenTypeBang)
		}
	case '=':
		if peek == '=' {
			s.current++
			s.addToken(start, TokenTypeEqualEqual)
		} else {
			s.addToken(start, TokenTypeEqual)
		}
	case '<':
		if peek == '=' {
			s.current++
			s.addToken(start, TokenTypeLessEqual)
		} else {
			s.addToken(start, TokenTypeLess)
		}
	case '>':
		if peek == '=' {
			s.current++
			s.addToken(start, TokenTypeGreaterEqual)
		} else {
			s.addToken(start, TokenTypeGreater)
		}

		// String handling
	case '"':
		for s.source[s.current] != '"' && s.current < len(s.source) {
			if s.source[s.current] == '\n' {
				s.line++
				s.visitedLinesLen = s.current
			}
			s.current++
		}
		if s.current >= len(s.source) {
			return fmt.Errorf("unterminated string (line %d pos %d)", s.line+1, s.current-s.visitedLinesLen)
		}
		s.current++
		s.addLiteralToken(start, TokenTypeString, string(s.source[start+1:s.current-1]))

	// Numbers
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		for s.current < len(s.source) && isDigit(s.source[s.current]) {
			s.current++
		}
		if s.source[s.current] == '.' && isDigit(s.source[s.current+1]) {
			// Advance to end of non-integer number
			s.current++
			for s.current < len(s.source) && isDigit(s.source[s.current]) {
				s.current++
			}
		}
		s.addLiteralToken(start, TokenTypeNumber, string(s.source[start:s.current]))

	default:
		if isAlpha(c) {
			for s.current < len(s.source) && isAlphaNumeric(s.source[s.current]) {
				s.current++
			}
			s.addIdentifier(start)
		} else {
			return fmt.Errorf("unexpected character: %c (line %d col %d)", c, s.line+1, s.current-s.visitedLinesLen)
		}
	}
	return nil
}

// addToken appends a non-literal token to the token list
func (s *Scanner) addToken(start int, tokenType TokenType) {
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Lexeme:  string(s.source[start:s.current]),
		Literal: nil,
		Pos: Pos{
			Line:  s.line,
			Start: start - s.visitedLinesLen,
			End:   s.current - s.visitedLinesLen,
		},
	})
}

func (s *Scanner) addLiteralToken(start int, tokenType TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Lexeme:  string(s.source[start:s.current]),
		Literal: literal,
		Pos: Pos{
			Line:  s.line,
			Start: start - s.visitedLinesLen,
			End:   s.current - s.visitedLinesLen,
		},
	})
}

func (s *Scanner) addIdentifier(start int) {
	if keyword, ok := ReservedKeywords[string(s.source[start:s.current])]; ok {
		s.addToken(start, keyword)
	} else {
		s.addToken(start, TokenTypeIdentifier)
	}
}
