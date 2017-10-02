package ast

import (
	"bytes"
	"fmt"
	"token"
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

type IntegerLiteral struct {
	Token    *token.Token
	IntValue int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.IntValue)
}

type PrefixExpression struct {
	Token      *token.Token
	Operator   string
	Expression Expression
}

func (pr *PrefixExpression) expressionNode() {}
func (pr *PrefixExpression) TokenLiteral() string {
	return pr.Token.Literal
}
func (pr *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pr.Operator, pr.Expression.String())
}

type InfixExpression struct {
	Token      *token.Token
	Left Expression
	Operator   string
	Right Expression
}

func (in *InfixExpression) expressionNode() {}
func (in *InfixExpression) TokenLiteral() string {
	return in.Token.Literal
}
func (in *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", in.Left.String(), in.Operator, in.Right.String())
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
	Token      *token.Token
	Expression Expression
}

func (s *ExpressionStatement) statementNode() {}
func (s *ExpressionStatement) TokenLiteral() string {
	return s.Token.Literal
}
func (s *ExpressionStatement) String() string {
	return s.Expression.String()
}
