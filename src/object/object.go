package object

import "fmt"

type Type string

const (
	TYPE_INTEGER = "INTEGER"
	TYPE_BOOLEAN = "BOOLEAN"
	TYPE_NULL    = "NULL"
	TYPE_RETURN  = "RETURN"
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
