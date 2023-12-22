package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	"github.com/aselhid/indoscript/internal/ast"
)

type Reader struct {
	*bufio.Reader
}

// TODO: handle error rune
func (r *Reader) PeekRune() rune {
	for peekBytes := 4; peekBytes > 0; peekBytes-- {
		b, err := r.Peek(peekBytes)
		if err != nil {
			continue
		}
		r, _ := utf8.DecodeRune(b)
		return r
	}
	return '\000'
}

type Scanner struct {
	reader     *Reader
	tokens     []ast.Token
	buffer     []rune
	lineNumber int
	stdErr     io.Writer
}

func NewScanner(r io.Reader, stdErr io.Writer) *Scanner {
	return &Scanner{
		reader:     &Reader{bufio.NewReader(r)},
		stdErr:     stdErr,
		lineNumber: 1,
	}
}

func (s *Scanner) ScanTokens() []ast.Token {
	for !s.isAtEnd() {
		s.scanToken()
		s.clearBuffer()
	}
	s.tokens = append(s.tokens, ast.Token{TokenType: ast.TokenEof, LineNumber: s.lineNumber})
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()
	switch char {
	case '(':
		s.addToken(ast.TokenLeftParenthesis)
	case ')':
		s.addToken(ast.TokenRightParenthesis)
	case '{':
		s.addToken(ast.TokenLeftBrace)
	case '}':
		s.addToken(ast.TokenRightBrace)
	case ',':
		s.addToken(ast.TokenComma)
	case '.':
		s.addToken(ast.TokenDot)
	case '+':
		s.addToken(ast.TokenPlus)
	case '-':
		s.addToken(ast.TokenMinus)
	case ';':
		s.addToken(ast.TokenSemicolon)
	case '*':
		s.addToken(ast.TokenStar)

	case '/':
		if s.match('/') {
			for !s.isAtEnd() {
				if r := s.reader.PeekRune(); r == '\n' {
					break
				}
				s.advance()
			}
		} else {
			s.addToken(ast.TokenSlash)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.TokenEqualEqual)
		} else {
			s.addToken(ast.TokenEqual)
		}
	case '!':
		if s.match('=') {
			s.addToken(ast.TokenBangEqual)
		} else {
			s.addToken(ast.TokenBang)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.TokenGreaterEqual)
		} else {
			s.addToken(ast.TokenGreater)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.TokenLessEqual)
		} else {
			s.addToken(ast.TokenLess)
		}

	// ignore white space
	case ' ':
	case '\t':
	case '\r':

	case '\n':
		s.lineNumber++

	case '"':
		s.buffer = []rune{char}
		s.string()

	default:
		if s.isDigit(char) {
			s.buffer = []rune{char}
			s.number()
		} else if s.isAllowedAlpha(char) {
			s.buffer = []rune{char}
			s.identifierOrKeyword()
		} else {
			s.error(fmt.Sprintf("found unexpected character  \"%c\"", char))
		}

	}
}

func (s *Scanner) string() {
	for s.reader.PeekRune() != '"' && s.reader.PeekRune() != '\n' {
		s.buffer = append(s.buffer, s.advance())
	}

	if s.isAtEnd() || s.reader.PeekRune() == '\n' {
		s.error("unterminated string")
		return
	}

	s.buffer = append(s.buffer, s.advance())
	value := string(s.buffer[1 : len(s.buffer)-1])
	s.addTokenWithLiteral(ast.TokenString, value)
}

func (s *Scanner) number() {
	for s.isDigit(s.reader.PeekRune()) {
		s.buffer = append(s.buffer, s.advance())
	}

	if r := s.reader.PeekRune(); r == '.' {
		s.buffer = append(s.buffer, s.advance())
		for s.isDigit(s.reader.PeekRune()) {
			s.buffer = append(s.buffer, s.advance())
		}
	}

	value, _ := strconv.ParseFloat(string(s.buffer), 64)
	s.addTokenWithLiteral(ast.TokenNumber, value)
}

func (s *Scanner) identifierOrKeyword() {
	for s.isAllowedAlphanumeric(s.reader.PeekRune()) {
		s.buffer = append(s.buffer, s.advance())
	}

	text := string(s.buffer)
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = ast.TokenIdentifier
	}
	s.addToken(tokenType)
}

// TODO: handle error
func (s *Scanner) advance() rune {
	char, _, _ := s.reader.ReadRune()
	return char
}

func (s *Scanner) match(target rune) bool {
	char, _, _ := s.reader.ReadRune()
	if target != char {
		s.reader.UnreadRune()
		return false
	}
	return true
}

func (s *Scanner) isAtEnd() bool {
	r := s.reader.PeekRune()
	return r == '\000'
}

func (s *Scanner) addToken(tokenType ast.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType ast.TokenType, literal any) {
	var lexeme string
	if len(s.buffer) > 0 {
		lexeme = string(s.buffer)
	}
	s.tokens = append(s.tokens, ast.Token{TokenType: tokenType, LineNumber: s.lineNumber, Literal: literal, Lexeme: lexeme})
}

func (s *Scanner) clearBuffer() {
	s.buffer = s.buffer[:0]
}

func (s *Scanner) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (s *Scanner) isAllowedAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}

func (s *Scanner) isAllowedAlphanumeric(r rune) bool {
	return s.isAllowedAlpha(r) || s.isDigit(r)
}

func (s *Scanner) error(message string) {
	s.stdErr.Write([]byte(fmt.Sprintf("[line %d] %s\n", s.lineNumber, message)))
}
