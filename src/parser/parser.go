package parser

import (
	"ast"
	"fmt"
	"lexer"
	"strconv"
	"token"
)

const (
	_        = iota
	LOWEST
	EQUALS
	INEQUALS
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       INEQUALS,
	token.GT:       INEQUALS,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

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
	res.prefixParseFns[token.INT] = res.parseIntegerLiteral
	res.prefixParseFns[token.BANG] = res.parsePrefixOperator
	res.prefixParseFns[token.MINUS] = res.parsePrefixOperator
	res.prefixParseFns[token.TRUE] = res.parseBooleanLiteral
	res.prefixParseFns[token.FALSE] = res.parseBooleanLiteral
	res.prefixParseFns[token.LPAREN] = res.parseGroupedExpression
	res.prefixParseFns[token.IF] = res.parseIfExpression
	res.prefixParseFns[token.FUNCTION] = res.parseFunctionExpression

	res.infixParseFns = make(map[token.Type]func(ast.Expression) ast.Expression)
	res.infixParseFns[token.EQ] = res.parseInfixExpression
	res.infixParseFns[token.NOT_EQ] = res.parseInfixExpression
	res.infixParseFns[token.LT] = res.parseInfixExpression
	res.infixParseFns[token.GT] = res.parseInfixExpression
	res.infixParseFns[token.PLUS] = res.parseInfixExpression
	res.infixParseFns[token.MINUS] = res.parseInfixExpression
	res.infixParseFns[token.ASTERISK] = res.parseInfixExpression
	res.infixParseFns[token.SLASH] = res.parseInfixExpression

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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	res := &ast.BlockStatement{Token: p.curToken}

	p.nextToken()

	statements := []ast.Statement{}
	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		s := p.parseStatement()
		if s != nil {
			statements = append(statements, s)
		}
		p.nextToken()
	}

	res.Statements = statements

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

func (p *Parser) parseIntegerLiteral() ast.Expression {
	number, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("cannot parse %q as integer", p.curToken.Literal))
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, IntValue: number}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	b := true
	if p.curToken.Literal == "false" {
		b = false
	}

	return &ast.BooleanLiteral{Token: p.curToken, BoolValue: b}
}

func (p *Parser) parsePrefixOperator() ast.Expression {
	cur := p.curToken

	p.nextToken()

	return &ast.PrefixExpression{
		Token:      cur,
		Operator:   cur.Literal,
		Expression: p.parseExpression(PREFIX),
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	res := &ast.InfixExpression{
		Left:     left,
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	res.Right = p.parseExpression(precedence)

	return res
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	tok := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	condition := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	consequence := p.parseBlockStatement()

	// check if there is an else
	var alternative *ast.BlockStatement
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		alternative = p.parseBlockStatement()
	}

	return &ast.IfExpression{
		Token:       tok,
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
	}
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	tok := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	params := []*ast.Identifier{}
	for {
		if p.curToken.Type == token.RPAREN {
			break
		}

		if p.curToken.Type != token.IDENT {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token at function parameter list: %q", p.curToken.Literal))
		}

		id := &ast.Identifier{Token: p.curToken, Name: p.curToken.Literal}
		if id != nil {
			params = append(params, id)
		}

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		} else if p.peekTokenIs(token.RPAREN) {
			// nop
		} else {
			p.errors = append(p.errors, fmt.Sprintf("unexpected token at function parameter list: %q", p.peekToken.Literal))
		}

		p.nextToken()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &ast.FunctionExpression{
		Token:  tok,
		Params: params,
		Body:   body,
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixParseFn := p.prefixParseFns[p.curToken.Type]
	if prefixParseFn == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parser function for %s", p.curToken.Type))
		return nil
	}

	leftExp := prefixParseFn()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixParseFns := p.infixParseFns[p.peekToken.Type]
		if infixParseFns == nil {
			p.errors = append(p.errors, fmt.Sprintf("no infix parser function for %s", p.curToken.Type))
			return leftExp
		}

		p.nextToken()
		leftExp = infixParseFns(leftExp)
	}

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

func (p *Parser) curPrecedence() int {
	if res, ok := precedences[p.curToken.Type]; ok {
		return res
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if res, ok := precedences[p.peekToken.Type]; ok {
		return res
	}
	return LOWEST
}
