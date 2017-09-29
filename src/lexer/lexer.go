package lexer

import "token"

type Lexer struct {
	input    string
	position int // points to the ch
	readPos  int
	ch       byte // current char
}

func New(input string) *Lexer {
	res := &Lexer{input: input}
	res.readChar()
	return res
}

func (lx *Lexer) NextToken() token.Token {
	var res token.Token

	newToken := func(typ token.TokenType, ch byte) token.Token {
		var literal string
		if ch == 0 {
			literal = ""
		} else {
			literal = string(ch)
		}
		return token.Token{Type: typ, Literal: literal}
	}

	switch lx.ch {
	case '=':
		res = newToken(token.ASSIGN, lx.ch)
	case ';':
		res = newToken(token.SEMICOLON, lx.ch)
	case '(':
		res = newToken(token.LPAREN, lx.ch)
	case ')':
		res = newToken(token.RPAREN, lx.ch)
	case '{':
		res = newToken(token.LBRACE, lx.ch)
	case '}':
		res = newToken(token.RBRACE, lx.ch)
	case ',':
		res = newToken(token.COMMA, lx.ch)
	case '+':
		res = newToken(token.PLUS, lx.ch)
	case 0:
		res = newToken(token.EOF, 0)
	}

	lx.readChar()
	return res
}

func (lx *Lexer) readChar() {
	if lx.readPos >= len(lx.input) {
		lx.ch = 0
	} else {
		lx.ch = lx.input[lx.readPos]
	}
	lx.position = lx.readPos
	lx.readPos++
}
