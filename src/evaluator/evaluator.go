package evaluator

import (
	"ast"
	"fmt"
	"object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.IntValue}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: node.BoolValue}
	}

	panic(fmt.Sprintf("unhandled case %T", node))
}

func evalStatements(ss []ast.Statement) object.Object {
	var res object.Object

	for _, s := range ss {
		res = Eval(s)
	}

	return res
}
