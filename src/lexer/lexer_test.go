package lexer

import (
	"testing"
	"token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for i, tt := range tests {
		tok := lexer.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test %d: token type is wrong. Expected %q, got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test %d: literal is wrong. Expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
