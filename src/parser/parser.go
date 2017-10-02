package parser

import (
	"lexer"
	"token"
	"ast"
	"fmt"
)

const (
	_        = iota
	LOWEST
	EQUALS
	INEQUALS
	SUM
	PRODUCAT
	PREFIX
	CALL
)

type Parser struct {
	lx     *lexer.Lexer
	errors []string

	curToken  *token.Token
	peekToken *token.Token

	prefixParseFns map[token.Type]func() ast.Expression
	infixParseFns  map[token.Type]func(ast.Expression) ast.Expression
}

func New(lx *lexer.Lexer) *Parser {
	res := &Parser{lx: lx}
	res.prefixParseFns = make(map[token.Type]func() ast.Expression)
	res.prefixParseFns[token.IDENT] = res.parseIdentifier

	// read two tokens so curToken and peekToken are set
	res.nextToken()
	res.nextToken()

	return res
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(typ token.Type) {
	p.errors = append(p.errors, fmt.Sprintf("next token is expected to be %q, got: %q", typ, p.peekToken.Type))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	tmp := p.lx.NextToken()
	p.peekToken = &tmp
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	res := &ast.ReturnStatement{Token: p.curToken}

	// TODO now we just skip until semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return res
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	res := &ast.ExpressionStatement{Token: p.curToken}

	res.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return res
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFn := p.prefixParseFns[p.curToken.Type]
	if prefixParseFn == nil {
		return nil
	}

	leftExp := prefixParseFn()
	return leftExp
}

// expectPeek checks if the peek token is of specified type,
// then advances to the next one if it is. If not it will record it as an error.
func (p *Parser) expectPeek(typ token.Type) bool {
	if p.peekTokenIs(typ) {
		p.nextToken()
		return true
	}
	p.peekError(typ)
	return false
}

func (p *Parser) peekTokenIs(typ token.Type) bool {
	return p.peekToken.Type == typ
}
