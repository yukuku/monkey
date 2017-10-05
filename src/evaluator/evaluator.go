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
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfix(node.Operator, left, right)
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

func evalInfix(operator string, left object.Object, right object.Object) object.Object {
	convertToInteger := func(obj object.Object) int64 {
		switch obj := obj.(type) {
		case *object.Boolean:
			if obj.Value {
				return 1
			} else {
				return 0
			}
		case *object.Null:
			return 0
		case *object.Integer:
			return obj.Value
		}
		panic(fmt.Sprintf("unhandled type for integer conversion %T", obj))
	}

	switch operator {
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "*":
		fallthrough
	case "/":
		leftint := convertToInteger(left)
		rightint := convertToInteger(right)

		switch operator {
		case "+":
			return &object.Integer{Value: leftint + rightint}
		case "-":
			return &object.Integer{Value: leftint - rightint}
		case "*":
			return &object.Integer{Value: leftint * rightint}
		case "/":
			return &object.Integer{Value: leftint / rightint}
		}
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
