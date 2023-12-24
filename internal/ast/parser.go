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

/*
Grammar (so far)
----------------
program         -> declaration* EOF
declaration     -> varDeclaration | statement | assignment | funcDeclaration
funcDeclaration -> "fungsi" function
function        -> IDENTIFIER "(" parameters? ")" block
parameters      -> IDENTIFIER ( "," IDENTIFIER )*
assignment      -> IDENTIFIER "=" expression ";"
varDeclaration  -> "misal" IDENTIFIER "=" expresion ";"
statement       -> exprStmt | printStmt | block | ifStmt | whileStmt | returnStmt
returnStmt      -> "return" expression? ";"
whileStmt       -> "untuk" expression "{" statement "}"
ifStmt          -> "jika" expression "{" statement "}" ( "lain" "{" statement "}" )?
block           -> "{" declaration* "}"
exprStmt        -> expression ";"
printStmt       -> "cetak" expression ";"
expression      -> logic_or
logic_or        -> logic_and ( "atau" logic_and )*
logic_and       -> equality ( "dan" equality )*
equality        -> comparison ( ("!=" | "==") comparison )*
comparison      -> term ( ( ">" | ">=" | "<" | "<=" ) term )*
term            -> factor ( ( "-" | "+" ) factor )*
factor          -> unary ( ( "/" | "*" ) unary )*
unary           -> ( "!" | "-" ) unary | call
call            -> primary ( "(" arguments? ")" )*
primary         -> FALSE | TRUE | NIL | NUMBER | STRING | group | IDENTIFIER
group           -> "(" expression ")"
arguments       -> expression ( "," expression )*
*/

type Parser struct {
	stdErr   io.Writer
	tokens   []Token
	current  int
	hasError bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// using Expr first
func (p *Parser) Parse() ([]Stmt, bool) {
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
	var result []Stmt
	for !p.isAtEnd() {
		result = append(result, p.declaration())
	}
	return result, p.hasError
}

func (p *Parser) declaration() Stmt {
	switch {
	case p.match(TokenFunction):
		return p.funcDeclaration()
	case p.match(TokenLet):
		return p.varDeclaration()
	case p.match(TokenIdentifier):
		if p.match(TokenEqual) {
			return p.assignment()
		} else {
			p.undoAdvance()
		}
	}
	return p.statement()
}

func (p *Parser) assignment() Stmt {
	identifier := p.previous()
	p.consume(TokenEqual, "expect '=' for declaration")

	value := p.expression()
	p.consume(TokenSemicolon, "expect ';' after statement")
	return NewVarStmt(identifier, value)
}

func (p *Parser) varDeclaration() Stmt {
	identifier := p.consume(TokenIdentifier, "expect variable name after 'mulai'")
	p.consume(TokenEqual, "identifier without initialization is not allowed")

	initializer := p.expression()
	p.consume(TokenSemicolon, "expect ';' after statement")

	return NewVarStmt(identifier, initializer)
}

func (p *Parser) funcDeclaration() Stmt {
	name := p.consume(TokenIdentifier, "expect identifier after fungsi declaration")
	p.consume(TokenLeftParenthesis, "expect opening '(' after fungsi declaration")
	var parameters []Token
	if p.peek().TokenType != TokenRightParenthesis {
		for {
			parameters = append(parameters, p.consume(TokenIdentifier, "expect parameter name"))
			if !p.match(TokenComma) {
				break
			}
		}
	}
	p.consume(TokenRightParenthesis, "expect closing ')' after fungsi declaration")
	p.consume(TokenLeftBrace, "expect opening '{' to define fungsi body")
	body := p.block()
	return NewFuncStmt(name, parameters, body)
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(TokenPrint):
		return p.printStmt()
	case p.match(TokenLeftBrace):
		return NewBlockStmt(p.block())
	case p.match(TokenIf):
		return p.ifStmt()
	case p.match(TokenLoop):
		return p.whileStmt()
	case p.match(TokenReturn):
		return p.returnStmt()
	}
	return p.exprStmt()
}

func (p *Parser) printStmt() Stmt {
	expr := p.expression()
	p.consume(TokenSemicolon, "expect ';' after statement")
	return NewPrintStmt(expr)
}

func (p *Parser) exprStmt() Stmt {
	expr := p.expression()
	p.consume(TokenSemicolon, "expect ';' after statement")
	return NewExprStmt(expr)
}

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for p.peek().TokenType != TokenRightBrace && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(TokenRightBrace, "expect '}' to close the scope")

	return statements
}

func (p *Parser) ifStmt() Stmt {
	condition := p.expression()

	p.consume(TokenLeftBrace, "expect block start after jika condition")
	thenStmt := p.block()
	if p.match(TokenElse) {
		p.consume(TokenLeftBrace, "expect block start '}' after lain statement")
		elseStmt := p.block()
		return NewIfStmt(condition, NewBlockStmt(thenStmt), NewBlockStmt(elseStmt))
	}
	return NewIfStmt(condition, NewBlockStmt(thenStmt), NewBlockStmt(nil))
}

func (p *Parser) whileStmt() Stmt {
	condition := p.expression()

	p.consume(TokenLeftBrace, "expect block start '{' after selama statement ")
	stmt := p.block()
	return NewWhileStmt(condition, NewBlockStmt(stmt))
}

func (p *Parser) returnStmt() Stmt {
	var value Expr
	if p.peek().TokenType != TokenSemicolon {
		value = p.expression()
	}
	p.consume(TokenSemicolon, "expect ';' after return statement")
	return NewReturnStmt(value)
}

func (p *Parser) expression() Expr {
	return p.logicOr()
}

func (p *Parser) logicOr() Expr {
	expr := p.logicAnd()
	for p.match(TokenOr) {
		operator := p.previous()
		right := p.logicAnd()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) logicAnd() Expr {
	expr := p.equality()
	for p.match(TokenAnd) {
		operator := p.previous()
		right := p.equality()
		expr = NewLogicalExpr(expr, operator, right)
	}
	return expr
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
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()

	for p.match(TokenLeftParenthesis) {
		expr = p.finishCall(expr)
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if p.peek().TokenType != TokenRightParenthesis {
		for {
			arguments = append(arguments, p.expression())
			if !p.match(TokenComma) {
				break
			}
		}
	}
	parenthesis := p.consume(TokenRightParenthesis, "expect closing ')' after fungsi call")
	return NewCallExpr(callee, arguments, parenthesis)
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
	case p.match(TokenIdentifier):
		return NewVarExpr(p.previous())
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

func (p *Parser) undoAdvance() {
	if p.current > 0 {
		p.current--
	}
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
		case TokenFunction, TokenLet, TokenLoop, TokenIf, TokenPrint, TokenReturn:
			return
		}
		p.advance()
	}
}

func (p *Parser) error(token Token, errMessage string) {
	fmt.Println(errMessage)
	location := "at the end of file"
	if token.TokenType != TokenEof {
		location = fmt.Sprintf("at '%s'", token.Lexeme)
	}
	err := parseError{message: fmt.Sprintf("[line %d] Error %s: %s\n", token.LineNumber, location, errMessage)}
	p.stdErr.Write([]byte(err.Error()))
	panic(err)
}
