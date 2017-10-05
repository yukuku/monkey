package evaluator

import (
	"ast"
	"fmt"
	"object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		// special case for Program, need to unwrap Return
		res := evalStatements(node.Statements)
		if r, ok := res.(*object.Return); ok {
			return r.Value
		} else {
			return res
		}
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
	case *ast.IfExpression:
		condition := Eval(node.Condition)

		if convertToBool(condition) {
			return Eval(node.Consequence)
		} else {
			if node.Alternative == nil {
				return &object.Null{}
			} else {
				return Eval(node.Alternative)
			}
		}
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.Value)}
	}

	panic(fmt.Sprintf("unhandled case %T", node))
}
func evalPrefix(operator string, operand object.Object) object.Object {
	switch operator {
	case "!":
		return &object.Boolean{Value: !convertToBool(operand)}
	case "-":
		return &object.Integer{Value: -convertToInteger(operand)}
	}

	panic(fmt.Sprintf("unhandled operator %s", operator))
}

func convertToInteger(obj object.Object) int64 {
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

func convertToBool(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	case *object.Integer:
		return obj.Value != 0
	}
	panic(fmt.Sprintf("unhandled type for bool conversion %T", obj))
}

func evalInfix(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	// these operators return integer
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

		// these operators return boolean
	case ">":
		fallthrough
	case "<":
		fallthrough
	case "==":
		fallthrough
	case "!=":
		leftint := convertToInteger(left)
		rightint := convertToInteger(right)

		switch operator {
		case ">":
			return &object.Boolean{Value: leftint > rightint}
		case "<":
			return &object.Boolean{Value: leftint < rightint}
		case "==":
			return &object.Boolean{Value: leftint == rightint}
		case "!=":
			return &object.Boolean{Value: leftint != rightint}
		}
	}

	panic(fmt.Sprintf("unhandled operator %s", operator))
}

func evalStatements(ss []ast.Statement) object.Object {
	var res object.Object

	for _, s := range ss {
		res = Eval(s)

		// if res is a Return, stop evaluating and return it immediately
		if r, ok := res.(*object.Return); ok {
			return r
		}
	}

	return res
}
