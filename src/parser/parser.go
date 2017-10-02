package parser

import (
	"lexer"
	"token"
	"ast"
)

type Parser struct {
	lx *lexer.Lexer

	curToken  *token.Token
	peekToken *token.Token
}

func New(lx *lexer.Lexer) *Parser {
	res := &Parser{lx: lx}

	// read two tokens so curToken and peekToken are set
	res.nextToken()
	res.nextToken()

	return res
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	tmp := p.lx.NextToken()
	p.peekToken = &tmp;
}

func (p *Parser) Parse() *ast.Program {
	res := &ast.Program{}

	for p.curToken.Type != token.EOF {
		s := p.parseStatement()
		if s != nil {
			res.Statements = append(res.Statements, s)
		}
		p.nextToken()
	}

	return res
}
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
	return nil
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	res := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	res.Ident = &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO now we just skip until semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return res
}

// expectPeek checks if the peek token is of specified type, then advances to the next
// one if it is.
func (p *Parser) expectPeek(typ token.Type) bool {
	if p.peekTokenIs(typ) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(typ token.Type) bool {
	return p.peekToken.Type == typ
}
