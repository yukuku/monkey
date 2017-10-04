package parser

import (
	"ast"
	"lexer"
	"testing"
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

	tests := []struct {
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

	if !testIdentifier(t, es.Expression, "foo") {
		return
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

	if !testIntegerLiteral(t, es.Expression, 7) {
		return
	}
}

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true; false;"
	lx := lexer.New(input)
	p := New(lx)
	prog := p.Parse()
	cannotHaveErrors(t, p)

	if len(prog.Statements) != 2 {
		t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
	}

	es, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("first statement is not an expression statement")
	}

	if !testBooleanLiteral(t, es.Expression, true) {
		return
	}

	es, ok = prog.Statements[1].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("second statement is not an expression statement")
	}

	if !testBooleanLiteral(t, es.Expression, false) {
		return
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input  string
		op     string
		intval int64
	}{
		{"!7;", "!", 7},
		{"-77;", "-", 77},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		prog := p.Parse()
		cannotHaveErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
		}

		es, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("first statement is not an expression statement")
		}

		pr, ok := es.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("first statement is not a prefix expression")
		}

		if pr.Operator != tt.op {
			t.Fatalf("operator is not correct. expected: %q, got: %q", tt.op, pr.Operator)
		}

		if !testIntegerLiteral(t, pr.Expression, tt.intval) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, intval int64) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp is not an integer literal. got: %T", exp)
		return false
	}

	if il.IntValue != intval {
		t.Errorf("int value is not correct. expected: %d, got: %d", intval, il.IntValue)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, b bool) bool {
	bl, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp is not a boolean literal. got: %T", exp)
		return false
	}

	if bl.BoolValue != b {
		t.Errorf("bool value is not correct. expected: %t, got: %t", b, bl.BoolValue)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, name string) bool {
	id, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not an identifier. got: %T", exp)
		return false
	}

	if id.Name != name {
		t.Errorf("identifier name is not correct. expected: %d, got: %d", name, id.Name)
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input string
		left  int64
		op    string
		right int64
	}{
		{"7 + 7", 7, "+", 7},
		{"7 - 7", 7, "-", 7},
		{"7 * 7", 7, "*", 7},
		{"7 / 7", 7, "/", 7},
		{"7 > 7", 7, ">", 7},
		{"7 < 7", 7, "<", 7},
		{"7 == 7", 7, "==", 7},
		{"7 != 7", 7, "!=", 7},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		prog := p.Parse()
		cannotHaveErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
		}

		es, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("first statement is not an expression statement")
		}

		in, ok := es.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("first statement is not a infix expression")
		}

		if in.Operator != tt.op {
			t.Fatalf("operator is not correct. expected: %q, got: %q", tt.op, in.Operator)
		}

		if !testIntegerLiteral(t, in.Left, tt.left) {
			return
		}

		if !testIntegerLiteral(t, in.Right, tt.right) {
			return
		}
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 < 5 == true == false",
			"(((3 < 5) == true) == false)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.in))
		prog := p.Parse()
		cannotHaveErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
		}

		if tt.out != prog.String() {
			t.Fatalf("wrong parsing. expected: %q, got: %q", tt.out, prog.String())
		}
	}
}

func TestIfExpressions(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		//{
		//	"if(x==1){y}",
		//	"if (x == 1) {y;}",
		//},
		//{
		//	"if(x==1){ y ; let z=3; }",
		//	"if (x == 1) {y; let z = 3;}",
		//},
		{
			"if(x<y+3){x;}else{y}",
			"if (x < (y + 3)) {x;} else {y;}",
		},
		//{
		//	"let z=if(x<y+3){x;}else{y}",
		//	"let z = if (x < y + 3) {x;} else {y;}",
		//},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.in))
		prog := p.Parse()
		cannotHaveErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
		}

		if tt.out != prog.String() {
			t.Fatalf("wrong parsing. expected: %q, got: %q", tt.out, prog.String())
		}
	}
}

func TestFunctionExpressions(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			"fn(){z}",
			"fn () {z;}",
		},
		{
			"fn(x,y){z;a}",
			"fn (x, y) {z;a;}",
		},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.in))
		prog := p.Parse()
		cannotHaveErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
		}

		if tt.out != prog.String() {
			t.Fatalf("wrong parsing. expected: %q, got: %q", tt.out, prog.String())
		}
	}
}

func TestCallExpressions(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			"tambah (1,3)",
			"tambah(1, 3)",
		},
		{
			"tambah(7,tambah(9,10),tambah())",
			"tambah(7, tambah(9, 10), tambah())",
		},
		{
			"fn(x){z}(y)(xxx,1,23)",
			"fn (x) {z;}(y)(xxx, 1, 23)",
		},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.in))
		prog := p.Parse()
		cannotHaveErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("wrong number of statements. got: %d", len(prog.Statements))
		}

		if tt.out != prog.String() {
			t.Fatalf("wrong parsing. expected: %q, got: %q", tt.out, prog.String())
		}
	}
}
