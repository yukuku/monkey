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
	case *ast.PrefixExpression:
		right := Eval(node.Expression)
		return evalPrefix(node.Operator, right)
	}

	panic(fmt.Sprintf("unhandled case %T", node))
}
func evalPrefix(operator string, operand object.Object) object.Object {
	switch operator {
	case "!":
		switch operand := operand.(type) {
		case *object.Boolean:
			if operand.Value {
				return &object.Boolean{Value: false}
			} else {
				return &object.Boolean{Value: true}
			}
		case *object.Null:
			return &object.Boolean{Value: true}
		case *object.Integer:
			intval := operand.Value
			if intval == 0 {
				return &object.Boolean{Value: true}
			} else {
				return &object.Boolean{Value: false}
			}
		}
		panic(fmt.Sprintf("unhandled bang operand %T", operand))
	case "-":
		switch operand := operand.(type) {
		case *object.Boolean:
			if operand.Value {
				return &object.Integer{Value: -1}
			} else {
				return &object.Integer{Value: 0}
			}
		case *object.Null:
			return &object.Integer{Value: 0}
		case *object.Integer:
			intval := operand.Value
			if intval == 0 {
				return &object.Integer{Value: -intval}
			} else {
				return &object.Integer{Value: -intval}
			}
		}
		panic(fmt.Sprintf("unhandled minus operand %T", operand))
	}

	panic(fmt.Sprintf("unhandled operator %s", operator))
}

func evalStatements(ss []ast.Statement) object.Object {
	var res object.Object

	for _, s := range ss {
		res = Eval(s)
	}

	return res
}
