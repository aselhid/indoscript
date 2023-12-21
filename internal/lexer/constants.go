package lexer

import "github.com/aselhid/indoscript/internal/ast"

var keywords = map[string]ast.TokenType{
	"misal":   ast.TokenLet,
	"jika":    ast.TokenIf,
	"lain":    ast.TokenElse,
	"fungsi":  ast.TokenFunction,
	"balikin": ast.TokenReturn,
	"kosong":  ast.TokenNil,
	"benar":   ast.TokenTrue,
	"salah":   ast.TokenFalse,
	"untuk":   ast.TokenFor,
	"cetak":   ast.TokenPrint,
	"dan":     ast.TokenAnd,
	"atau":    ast.TokenOr,
}
