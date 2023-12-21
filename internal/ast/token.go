package ast

type TokenType uint8

const (
	// Keywords
	TokenLet      TokenType = iota // misal
	TokenIf                        // jika
	TokenElse                      // lain
	TokenFunction                  // fungsi
	TokenReturn                    // balikin
	TokenNil                       // kosong
	TokenTrue                      // benar
	TokenFalse                     // salah
	TokenFor                       // untuk
	TokenPrint                     // cetak -- TODO: make as part of std library
	TokenAnd                       // dan
	TokenOr                        // atau

	// Single character token
	TokenLeftParenthesis  // (
	TokenRightParenthesis // )
	TokenLeftBrace        // {
	TokenRightBrace       // }
	TokenComma            // ,
	TokenDot              // .
	TokenPlus             // +
	TokenMinus            // -
	TokenSemicolon        // ;
	TokenStar             // *

	// Single & double characters token
	TokenSlash        // /
	TokenEqual        // =
	TokenEqualEqual   // ==
	TokenBang         // !
	TokenBangEqual    // !=
	TokenGreater      // >
	TokenGreaterEqual // >=
	TokenLess         // <
	TokenLessEqual    // <=

	// Literals
	TokenIdentifier
	TokenString
	TokenNumber

	TokenEof
)

type Token struct {
	TokenType  TokenType
	LineNumber int
	Lexeme     string
	Literal    any
}
