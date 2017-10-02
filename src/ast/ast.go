package ast

import (
	"bytes"
	"token"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		var buf bytes.Buffer
		for _, s := range p.Statements {
			buf.WriteString(s.TokenLiteral())
			buf.WriteString(" ")
		}
		return buf.String()
	} else {
		return ""
	}
}

type Identifier struct {
	Token *token.Token
	Name  string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type LetStatement struct {
	Token *token.Token
	Ident *Identifier
	Value Expression
}

func (s *LetStatement) statementNode() {}
func (s *LetStatement) TokenLiteral() string {
	return s.Token.Literal
}
