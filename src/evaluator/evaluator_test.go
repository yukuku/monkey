package evaluator

import (
	"lexer"
	"object"
	"parser"
	"testing"
)

func testEval(in string) object.Object {
	lx := lexer.New(in)
	p := parser.New(lx)
	prog := p.Parse()
	return Eval(prog)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"7", 7},
		{"777", 777},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}

func testIntegerObject(t *testing.T, o object.Object, expected int64) bool {
	i, ok := o.(*object.Integer)
	if !ok {
		t.Errorf("object is not integer. got: %T", o)
		return false
	}

	if i.Value != expected {
		t.Errorf("integer is expected to be %d, got %d", expected, i.Value)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		if !testBooleanObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}

func testBooleanObject(t *testing.T, o object.Object, expected bool) bool {
	b, ok := o.(*object.Boolean)
	if !ok {
		t.Errorf("object is not boolean. got: %T", o)
		return false
	}

	if b.Value != expected {
		t.Errorf("boolean is expected to be %t, got %t", expected, b.Value)
		return false
	}

	return true
}

func TestEvalBang(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!7", false},
		{"!!7", true},
	}

	for _, tt := range tests {
		if !testBooleanObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}

func TestEvalMinus(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"-5", -5},
		{"--7", 7},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}

func TestEvalIntegerExpressions(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}

	for _, tt := range tests {
		if !testBooleanObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}
