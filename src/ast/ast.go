package ast

import (
	"bytes"
	"token"
	"fmt"
)

type Node interface {
	TokenLiteral() string
	String() string
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
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	if len(p.Statements) > 0 {
		buf := bytes.Buffer{}
		for _, s := range p.Statements {
			buf.WriteString(s.String())
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
func (i *Identifier) String() string {
	return i.Name
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
func (s *LetStatement) String() string {
	return fmt.Sprintf("%s %s = %s;", s.TokenLiteral(), s.Ident.Name, s.Value.String())
}

type ReturnStatement struct {
	Token *token.Token
	Value Expression
}

func (s *ReturnStatement) statementNode() {}
func (s *ReturnStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ReturnStatement) String() string {
	return fmt.Sprintf("%s %s;", s.TokenLiteral(), s.Value.String())
}

type ExpressionStatement struct {
	Token *token.Token
	Expression Expression
}
func (s *ExpressionStatement) statementNode() {}
func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ExpressionStatement) String() string {
	return s.Expression.String()
}
