package evaluator

import (
	"testing"
	"object"
	"lexer"
	"parser"
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
