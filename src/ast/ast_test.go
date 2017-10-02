package ast

import (
	"testing"
	"token"
)

func TestString(t *testing.T) {
	prog := Program{
		Statements: []Statement{
			&LetStatement{
				Token: &token.Token{Type:token.LET, Literal:"let"},
				Ident: &Identifier{
					Token: &token.Token{Type:token.IDENT, Literal:"satu"},
					Name:"satu",
				},
				Value: &Identifier{
					Token: &token.Token{Type:token.IDENT, Literal:"dua"},
					Name:"dua",
				},
			},
		},
	}

	if prog.String() != "let satu = dua;" {
		t.Fatalf("String() in program is wrong. got: %s", prog.String())
	}
}
