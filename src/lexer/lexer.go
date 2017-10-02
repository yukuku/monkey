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

	newToken := func(typ token.Type, ch byte) token.Token {
		var literal string
		if ch == 0 {
			literal = ""
		} else {
			literal = string(ch)
		}
		return token.Token{Type: typ, Literal: literal}
	}

	// skip whitespaces
	for lx.ch == ' ' || lx.ch == '\t' || lx.ch == '\n' || lx.ch == '\r' {
		lx.readChar()
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
	case '-':
		res = newToken(token.MINUS, lx.ch)
	case '!':
		res = newToken(token.BANG, lx.ch)
	case '/':
		res = newToken(token.SLASH, lx.ch)
	case '*':
		res = newToken(token.ASTERISK, lx.ch)
	case '<':
		res = newToken(token.LT, lx.ch)
	case '>':
		res = newToken(token.GT, lx.ch)
	case 0:
		res = newToken(token.EOF, 0)
	default:
		isLetter := func(ch byte) bool {
			return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
		}

		isDigit := func(ch byte) bool {
			return '0' <= ch && ch <= '9'
		}

		if isLetter(lx.ch) {
			readIdentifier := func() string {
				pos := lx.position
				for isLetter(lx.ch) {
					lx.readChar()
				}
				return lx.input[pos:lx.position]
			}
			res.Literal = readIdentifier()
			res.Type = token.LookupIdent(res.Literal)
			return res

		} else if isDigit(lx.ch) {
			readNumber := func() string {
				pos := lx.position
				for isDigit(lx.ch) {
					lx.readChar()
				}
				return lx.input[pos:lx.position]
			}
			res.Type = token.INT
			res.Literal = readNumber()
			return res

		} else {
			res = newToken(token.ILLEGAL, lx.ch)
		}
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
