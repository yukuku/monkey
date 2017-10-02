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

	cannotHaveErrors(t, p)

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

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10 + 5;
return add(4, 5);
 	`

	lx := lexer.New(input)
	p := New(lx)

	prog := p.Parse()
	if prog == nil {
		t.Fatalf("Parse() returned nil")
	}

	cannotHaveErrors(t, p)

	if len(prog.Statements) != 3 {
		t.Fatalf("Statement count is wrong: %d", len(prog.Statements))
	}

	for _, s := range prog.Statements {
		rs, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is not return statement. got: %T", s)
		}
		if rs.TokenLiteral() != "return" {
			t.Errorf("statement token must be 'return'. got: %q", rs.TokenLiteral())
		}
	}
}

func cannotHaveErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) != 0 {
		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("- error: %s", msg)
		}
		t.FailNow()
	}
}

func TestErrorReporting(t *testing.T) {
	input := `
let x = 5;
let y 10;
let 129123;
 	`

	lx := lexer.New(input)
	p := New(lx)

	prog := p.Parse()
	if prog == nil {
		t.Fatalf("Parse() returned nil")
	}

	// must have errors
	errors := p.Errors()
	expectedErrorCount := 2
	if len(errors) != expectedErrorCount {
		t.Fatalf("parser is expected to have %d errors. got: %d", expectedErrorCount, len(errors))
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foo;"
	lx := lexer.New(input)
	p := New(lx)
	prog := p.Parse()
	cannotHaveErrors(t, p)

	if len(prog.Statements) != 1 {
		t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
	}

	es, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("first statement is not an expression statement")
	}

	id, ok := es.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("first statement is not an identifier")
	}

	if id.Name != "foo" {
		t.Fatalf("name of identifier is not correct. got: %s", id.Name)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "7;"
	lx := lexer.New(input)
	p := New(lx)
	prog := p.Parse()
	cannotHaveErrors(t, p)

	if len(prog.Statements) != 1 {
		t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
	}

	es, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("first statement is not an expression statement")
	}

	il, ok := es.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("first statement is not an integer literal")
	}

	if il.IntValue != 7 {
		t.Fatalf("name of identifier is not correct. got: %d", il.IntValue)
	}
}
