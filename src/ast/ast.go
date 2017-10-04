package ast

import (
	"bytes"
	"fmt"
	"strings"
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

type BooleanLiteral struct {
	Token     *token.Token
	BoolValue bool
}

func (bl *BooleanLiteral) expressionNode() {}
func (bl *BooleanLiteral) TokenLiteral() string {
	return bl.Token.Literal
}
func (bl *BooleanLiteral) String() string {
	if bl.BoolValue {
		return "true"
	} else {
		return "false"
	}
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
	Token    *token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (in *InfixExpression) expressionNode() {}
func (in *InfixExpression) TokenLiteral() string {
	return in.Token.Literal
}
func (in *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", in.Left.String(), in.Operator, in.Right.String())
}

type IfExpression struct {
	Token       *token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *IfExpression) String() string {
	if ie.Alternative == nil {
		return fmt.Sprintf("%s %s %s", ie.TokenLiteral(), ie.Condition, ie.Consequence)
	}
	return fmt.Sprintf("%s %s %s else %s", ie.TokenLiteral(), ie.Condition, ie.Consequence, ie.Alternative)
}

type FunctionExpression struct {
	Token  *token.Token
	Params []*Identifier
	Body   *BlockStatement
}

func (fu *FunctionExpression) expressionNode() {}
func (fu *FunctionExpression) TokenLiteral() string {
	return fu.Token.Literal
}
func (fu *FunctionExpression) String() string {
	names := []string{}
	for _, par := range fu.Params {
		names = append(names, par.Name)
	}
	paramsString := strings.Join(names, ", ")

	return fmt.Sprintf("%s (%s) %s", fu.TokenLiteral(), paramsString, fu.Body)
}

type CallExpression struct {
	Token     *token.Token
	Function  Expression
	Arguments []Expression
}

func (ca *CallExpression) expressionNode() {}
func (ca *CallExpression) TokenLiteral() string {
	return ca.Token.Literal
}
func (ca *CallExpression) String() string {
	args := []string{}
	for _, arg := range ca.Arguments {
		args = append(args, arg.String())
	}
	argsString := strings.Join(args, ", ")

	return fmt.Sprintf("%s(%s)", ca.Function, argsString)
}

type BlockStatement struct {
	Token      *token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	buf := bytes.Buffer{}
	buf.WriteString("{")
	for _, s := range bs.Statements {
		buf.WriteString(s.String())
		buf.WriteString(";")
	}
	buf.WriteString("}")
	return buf.String()
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
