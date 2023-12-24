package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aselhid/indoscript/internal/ast"
	"github.com/google/go-cmp/cmp"
)

func TestParentheses(t *testing.T) {
	scanner, stdErr := setupScanner("()\n)(")

	expected := []ast.Token{
		{TokenType: ast.TokenLeftParenthesis, LineNumber: 1},
		{TokenType: ast.TokenRightParenthesis, LineNumber: 1},
		{TokenType: ast.TokenRightParenthesis, LineNumber: 2},
		{TokenType: ast.TokenLeftParenthesis, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestBraces(t *testing.T) {
	scanner, stdErr := setupScanner("{}\n}{")

	expected := []ast.Token{
		{TokenType: ast.TokenLeftBrace, LineNumber: 1},
		{TokenType: ast.TokenRightBrace, LineNumber: 1},
		{TokenType: ast.TokenRightBrace, LineNumber: 2},
		{TokenType: ast.TokenLeftBrace, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestComma(t *testing.T) {
	scanner, stdErr := setupScanner(",\n,")

	expected := []ast.Token{
		{TokenType: ast.TokenComma, LineNumber: 1},
		{TokenType: ast.TokenComma, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestDot(t *testing.T) {
	scanner, stdErr := setupScanner(".\n.")

	expected := []ast.Token{
		{TokenType: ast.TokenDot, LineNumber: 1},
		{TokenType: ast.TokenDot, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestPlus(t *testing.T) {
	scanner, stdErr := setupScanner("+\n+")

	expected := []ast.Token{
		{TokenType: ast.TokenPlus, LineNumber: 1},
		{TokenType: ast.TokenPlus, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestMinus(t *testing.T) {
	scanner, stdErr := setupScanner("-\n-")

	expected := []ast.Token{
		{TokenType: ast.TokenMinus, LineNumber: 1},
		{TokenType: ast.TokenMinus, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestSemicolon(t *testing.T) {
	scanner, stdErr := setupScanner(";\n;")

	expected := []ast.Token{
		{TokenType: ast.TokenSemicolon, LineNumber: 1},
		{TokenType: ast.TokenSemicolon, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestStar(t *testing.T) {
	scanner, stdErr := setupScanner("*\n*")

	expected := []ast.Token{
		{TokenType: ast.TokenStar, LineNumber: 1},
		{TokenType: ast.TokenStar, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestSlash(t *testing.T) {
	scanner, stdErr := setupScanner("/\n///\n/")

	expected := []ast.Token{
		{TokenType: ast.TokenSlash, LineNumber: 1},
		{TokenType: ast.TokenSlash, LineNumber: 3},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestEqual(t *testing.T) {
	scanner, stdErr := setupScanner("=\n===")

	expected := []ast.Token{
		{TokenType: ast.TokenEqual, LineNumber: 1},
		{TokenType: ast.TokenEqualEqual, LineNumber: 2},
		{TokenType: ast.TokenEqual, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestBang(t *testing.T) {
	scanner, stdErr := setupScanner("!\n!!=")

	expected := []ast.Token{
		{TokenType: ast.TokenBang, LineNumber: 1},
		{TokenType: ast.TokenBang, LineNumber: 2},
		{TokenType: ast.TokenBangEqual, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestGreater(t *testing.T) {
	scanner, stdErr := setupScanner(">\n>>=")

	expected := []ast.Token{
		{TokenType: ast.TokenGreater, LineNumber: 1},
		{TokenType: ast.TokenGreater, LineNumber: 2},
		{TokenType: ast.TokenGreaterEqual, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestLess(t *testing.T) {
	scanner, stdErr := setupScanner("<\n<<=")

	expected := []ast.Token{
		{TokenType: ast.TokenLess, LineNumber: 1},
		{TokenType: ast.TokenLess, LineNumber: 2},
		{TokenType: ast.TokenLessEqual, LineNumber: 2},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestWhitespaces(t *testing.T) {
	scanner, stdErr := setupScanner("\r \t\n\n!")

	expected := []ast.Token{
		{TokenType: ast.TokenBang, LineNumber: 3},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestTerminatedString(t *testing.T) {
	scanner, stdErr := setupScanner("\"hello\"\n\"world\"")

	expected := []ast.Token{
		{TokenType: ast.TokenString, LineNumber: 1, Lexeme: "\"hello\"", Literal: "hello"},
		{TokenType: ast.TokenString, LineNumber: 2, Lexeme: "\"world\"", Literal: "world"},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestUnterminatedString(t *testing.T) {
	scanner, stdErr := setupScanner("\"hello\"\n\"world\n!")

	expected := []ast.Token{
		{TokenType: ast.TokenString, LineNumber: 1, Lexeme: "\"hello\"", Literal: "hello"},
		{TokenType: ast.TokenBang, LineNumber: 3},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	if stdErr.Len() == 0 || !strings.Contains(stdErr.String(), "unterminated string") {
		t.Fatal("expected unterminated string error in stdErr, found nothing")
	}
}

func TestNumber(t *testing.T) {
	scanner, stdErr := setupScanner("1234.0\n.0123\n0.1 2")

	expected := []ast.Token{
		{TokenType: ast.TokenNumber, LineNumber: 1, Lexeme: "1234.0", Literal: float64(1234.0)},
		{TokenType: ast.TokenDot, LineNumber: 2},
		{TokenType: ast.TokenNumber, LineNumber: 2, Lexeme: "0123", Literal: float64(123)},
		{TokenType: ast.TokenNumber, LineNumber: 3, Lexeme: "0.1", Literal: float64(0.1)},
		{TokenType: ast.TokenNumber, LineNumber: 3, Lexeme: "2", Literal: float64(2)},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestKeyword(t *testing.T) {
	scanner, stdErr := setupScanner("misal \njika lain fungsi balikin kosong benar salah selama cetak dan atau")
	expected := []ast.Token{
		{TokenType: ast.TokenLet, LineNumber: 1, Lexeme: "misal"},
		{TokenType: ast.TokenIf, LineNumber: 2, Lexeme: "jika"},
		{TokenType: ast.TokenElse, LineNumber: 2, Lexeme: "lain"},
		{TokenType: ast.TokenFunction, LineNumber: 2, Lexeme: "fungsi"},
		{TokenType: ast.TokenReturn, LineNumber: 2, Lexeme: "balikin"},
		{TokenType: ast.TokenNil, LineNumber: 2, Lexeme: "kosong"},
		{TokenType: ast.TokenTrue, LineNumber: 2, Lexeme: "benar"},
		{TokenType: ast.TokenFalse, LineNumber: 2, Lexeme: "salah"},
		{TokenType: ast.TokenLoop, LineNumber: 2, Lexeme: "selama"},
		{TokenType: ast.TokenPrint, LineNumber: 2, Lexeme: "cetak"},
		{TokenType: ast.TokenAnd, LineNumber: 2, Lexeme: "dan"},
		{TokenType: ast.TokenOr, LineNumber: 2, Lexeme: "atau"},
	}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	checkStdErrEmpty(t, stdErr)
}

func TestIdentifier(t *testing.T) {
	{
		scanner, stdErr := setupScanner("mimisal \njikaka lalainin fffungsi balikinnn kokokosong benarbenar sasalahlah se lama ce\ntak dandan watau")
		expected := []ast.Token{
			{TokenType: ast.TokenIdentifier, LineNumber: 1, Lexeme: "mimisal"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "jikaka"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "lalainin"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "fffungsi"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "balikinnn"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "kokokosong"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "benarbenar"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "sasalahlah"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "un"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "se"},
			{TokenType: ast.TokenIdentifier, LineNumber: 2, Lexeme: "lama"},
			{TokenType: ast.TokenIdentifier, LineNumber: 3, Lexeme: "tak"},
			{TokenType: ast.TokenIdentifier, LineNumber: 3, Lexeme: "dandan"},
			{TokenType: ast.TokenIdentifier, LineNumber: 3, Lexeme: "watau"},
		}
		actual := scanner.ScanTokens()
		compareTokens(t, expected, actual)
		checkStdErrEmpty(t, stdErr)
	}
}

func TestUnexpectedCharacter(t *testing.T) {
	scanner, stdErr := setupScanner("ඞ")
	expected := []ast.Token{}
	actual := scanner.ScanTokens()
	compareTokens(t, expected, actual)
	if stdErr.Len() == 0 || !strings.Contains(stdErr.String(), "found unexpected character") {
		t.Fatal("expected unexpected character error, found nothing")
	}
}

func compareTokens(t *testing.T, expected, actual []ast.Token) {
	if len(expected)+1 != len(actual) || actual[len(actual)-1].TokenType != ast.TokenEof {
		fmt.Printf("expected has %d elements while actual has %d elements\n", len(expected), len(actual))
		fmt.Printf("expected: %+v\n", expected)
		fmt.Printf("actual: %+v\n", actual)
		t.FailNow()
	}

	for i, expectedToken := range expected {
		if !cmp.Equal(expectedToken, actual[i]) {
			t.Fatalf("expected is %#v while actual is %#v", expectedToken, actual[i])
		}
	}
}

func checkStdErrEmpty(t *testing.T, stdErr *strings.Builder) {
	if stdErr.Len() != 0 {
		t.Fatalf("stdErr is not empty, %s", stdErr.String())
	}
}

func setupScanner(testcase string) (*Scanner, *strings.Builder) {
	r := strings.NewReader(testcase)
	w := new(strings.Builder)
	return NewScanner(r, w), w
}
