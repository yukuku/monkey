package parser

import (
	"testing"
	"lexer"
	"ast"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foo = 129123;
 	`

	lx := lexer.New(input)
	p := New(lx)

	prog := p.Parse()
	if prog == nil {
		t.Fatalf("Parse() returned nil")
	}

	tests := [] struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	if len(prog.Statements) != len(tests) {
		t.Fatalf("Statement count is wrong: %d", len(prog.Statements))
	}

	testLetStatement := func(s ast.Statement, identName string) bool {
		if s.TokenLiteral() != "let" {
			t.Errorf("statement token must be 'let'. got: %q", s.TokenLiteral())
			return false
		}

		letS, ok := s.(*ast.LetStatement)
		if !ok {
			t.Errorf("statement is not a LetStatement. got: %T", s)
			return false
		}

		if letS.Ident.Name != identName {
			t.Errorf("identifier name is not %q. got: %q", identName, letS.Ident.Name)
			return false
		}

		if letS.Ident.TokenLiteral() != identName {
			t.Errorf("identifier token literal is not %q. got: %q", identName, letS.Ident.TokenLiteral())
			return false
		}

		return true
	}

	for i, tt := range tests {
		s := prog.Statements[i]
		if !testLetStatement(s, tt.expectedIdentifier) {
			return
		}
	}
}
