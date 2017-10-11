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
	env := object.NewEnvironment()
	return Eval(prog, env)
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

func TestEvalIfExpressions(t *testing.T) {
	tests := []struct {
		in  string
		out interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { true }", true},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		ev := testEval(tt.in)
		if b, ok := tt.out.(bool); ok {
			if !testBooleanObject(t, ev, b) {
				return
			}
		} else if i, ok := tt.out.(int); ok {
			if !testIntegerObject(t, ev, int64(i)) {
				return
			}
		} else {
			if !testNullObject(t, ev) {
				return
			}
		}
	}
}

func testNullObject(t *testing.T, o object.Object) bool {
	_, ok := o.(*object.Null)
	if !ok {
		t.Errorf("object is not null. got: %T", o)
		return false
	}

	return true
}

func TestEvalReturn(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if (true) { if (true) { return 10; } return 1; }`, 10},
		{`if (true) { if (false) { return 10; } return 1; }`, 1},
		{`if (false) { if (true) { return 10; } return 1; } 25`, 25},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		in  string
		msg string
	}{
		{"return false + 3;", "first operand of + cannot be boolean"},
		{"false - 3;", "first operand of - cannot be boolean"},
		{"3 * false;", "second operand of * cannot be boolean"},
		{"false / true;", "first operand of / cannot be boolean"},
		{"false < true", "first operand of < cannot be boolean"},
		{"3 > true", "second operand of > cannot be boolean"},
		{"if (true) { if (true) { return 3 > true } }", "second operand of > cannot be boolean"},
		{"(1 != false) * 2", "cannot do != of different types"},
		{"1 * (3 == true)", "cannot do == of different types"},
		{"foo", "unknown identifier: foo"},
	}

	for _, tt := range tests {
		eval := testEval(tt.in)
		if e, ok := eval.(*object.Error); !ok {
			t.Errorf("result is not error. got: %v", eval)
		} else {
			if e.Message != tt.msg {
				t.Errorf("error message is wrong. expected: %q, got: %q", tt.msg, e.Message)
			}
		}
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		if !testIntegerObject(t, testEval(tt.in), tt.out) {
			return
		}
	}
}
