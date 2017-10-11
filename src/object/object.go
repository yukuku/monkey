package object

import (
	"ast"
	"fmt"
	"strings"
)

type Type int

const (
	TYPE_INTEGER  = iota + 1
	TYPE_BOOLEAN
	TYPE_NULL
	TYPE_RETURN
	TYPE_ERROR
	TYPE_FUNCTION
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() Type {
	return TYPE_INTEGER
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
func (b *Boolean) Type() Type {
	return TYPE_BOOLEAN
}

type Null struct {
}

func (n *Null) Inspect() string {
	return "null"
}
func (n *Null) Type() Type {
	return TYPE_NULL
}

type Return struct {
	Value Object
}

func (r *Return) Inspect() string {
	return fmt.Sprintf("return %s", r.Value.Inspect())
}
func (r *Return) Type() Type {
	return TYPE_RETURN
}

type Error struct {
	Message string
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR(%q)", e.Message)
}
func (e *Error) Type() Type {
	return TYPE_ERROR
}

type Function struct {
	Params []string
	Body   *ast.BlockStatement
	Env    *Environment
}

func (f *Function) Inspect() string {
	return fmt.Sprintf("fn (%s) %s", strings.Join(f.Params, ", "), f.Body)
}
func (f *Function) Type() Type {
	return TYPE_FUNCTION
}
