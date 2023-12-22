package ast

import (
	"fmt"
	"io"
)

type parseError struct {
	message string
}

func (e parseError) Error() string {
	return e.message
}

type Parser struct {
	tokens   []Token
	current  int
	stdErr   io.Writer
	hasError bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// using Expr first
func (p *Parser) Parse() ([]Expr, bool) {
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(parseError); ok {
				p.hasError = true
				p.sync()
			} else {
				panic(err)
			}
		}
	}()
	var result []Expr
	for !p.isAtEnd() {
		result = append(result, p.expression())
	}
	return result, p.hasError
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(TokenBangEqual, TokenEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		operator := p.previous()
		right := p.term()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(TokenMinus, TokenPlus) {
		operator := p.previous()
		right := p.factor()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(TokenSlash, TokenStar) {
		operator := p.previous()
		right := p.unary()
		expr = NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenBang, TokenMinus) {
		operator := p.previous()
		right := p.unary()
		return NewUnaryExpr(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(TokenFalse):
		return NewPrimaryExpr(false)
	case p.match(TokenTrue):
		return NewPrimaryExpr(true)
	case p.match(TokenNil):
		return NewPrimaryExpr(nil)
	case p.match(TokenNumber, TokenString):
		return NewPrimaryExpr(p.previous().Literal)
	case p.match(TokenLeftParenthesis):
		expr := p.expression()
		p.consume(TokenRightParenthesis, "Expect ')' after using '(' to group expression.")
		return expr
	}
	p.error(p.peek(), fmt.Sprintf("expecting expression, got %+v", p.previous()))
	return nil
}

func (p *Parser) consume(tokenType TokenType, errMessage string) Token {
	if p.peek().TokenType == tokenType {
		return p.advance()
	}
	p.error(p.peek(), errMessage)
	return Token{}
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if tokenType == p.peek().TokenType {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.current].TokenType == TokenEof
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) sync() {
	p.advance()

	for p.isAtEnd() {
		if p.previous().TokenType == TokenSemicolon {
			return
		}

		switch p.peek().TokenType {
		case TokenFunction, TokenLet, TokenFor, TokenIf, TokenPrint, TokenReturn:
			return
		}
		p.advance()
	}
}

func (p *Parser) error(token Token, errMessage string) {
	location := "at the end of file"
	if token.TokenType != TokenEof {
		location = fmt.Sprintf("at '%s'", token.Lexeme)
	}
	err := parseError{message: fmt.Sprintf("[line %d] Error %s: %s\n", token.LineNumber, location, errMessage)}
	p.stdErr.Write([]byte(err.Error()))
	panic(err)
}
