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

		pred, err := convertToBool(condition)
		if err != nil {
			return err
		}

		if pred {
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
	default:
		return newError("unhandled case %T", node)
	}
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func evalPrefix(operator string, operand object.Object) object.Object {
	if _, ok := operand.(*object.Error); ok {
		return operand
	}

	switch operator {
	case "!":
		val, err := convertToBool(operand)
		if err != nil {
			return err
		}
		return &object.Boolean{Value: !val}
	case "-":
		val, err := convertToInteger(operand)
		if err != nil {
			return err
		}
		return &object.Integer{Value: -val}
	}

	return newError("unhandled operator %s", operator)
}

func convertToInteger(obj object.Object) (result int64, err *object.Error) {
	switch obj := obj.(type) {
	case *object.Boolean:
		if obj.Value {
			result = 1
		} else {
			result = 0
		}
	case *object.Null:
		result = 0
	case *object.Integer:
		result = obj.Value
	default:
		err = newError("unhandled type for integer conversion %T", obj)
	}
	return
}

func convertToBool(obj object.Object) (result bool, err *object.Error) {
	switch obj := obj.(type) {
	case *object.Boolean:
		result = obj.Value
	case *object.Null:
		result = false
	case *object.Integer:
		result = obj.Value != 0
	default:
		err = newError("unhandled type for bool conversion %T", obj)
	}
	return
}

func evalInfix(operator string, left object.Object, right object.Object) object.Object {
	if _, ok := left.(*object.Error); ok {
		return left
	}
	if _, ok := right.(*object.Error); ok {
		return right
	}

	switch operator {
	// these operators return integer
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "*":
		fallthrough
	case "/":
		// booleans are not allowed
		if _, ok := left.(*object.Boolean); ok {
			return newError("first operand of %s cannot be boolean", operator)
		}
		if _, ok := right.(*object.Boolean); ok {
			return newError("second operand of %s cannot be boolean", operator)
		}

		leftint, err := convertToInteger(left)
		if err != nil {
			return err
		}

		rightint, err := convertToInteger(right)
		if err != nil {
			return err
		}

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

		// these operators return boolean:
	case ">":
		fallthrough
	case "<":
		fallthrough
	case "==":
		fallthrough
	case "!=":
		if operator == "<" || operator == ">" {
			// booleans are not allowed
			if _, ok := left.(*object.Boolean); ok {
				return newError("first operand of %s cannot be boolean", operator)
			}
			if _, ok := right.(*object.Boolean); ok {
				return newError("second operand of %s cannot be boolean", operator)
			}
		}
		if operator == "==" || operator == "!=" {
			// must be both boolean or both integers
			_, leftInt := left.(*object.Integer)
			_, rightInt := right.(*object.Integer)
			_, leftBool := left.(*object.Boolean)
			_, rightBool := right.(*object.Boolean)
			if (leftInt && rightInt) || (leftBool && rightBool) {
			} else {
				return newError("cannot do %s of different types", operator)
			}
		}

		leftint, err := convertToInteger(left)
		if err != nil {
			return err
		}
		rightint, err := convertToInteger(right)
		if err != nil {
			return err
		}

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

	return newError("unhandled operator %s", operator)
}

func evalStatements(ss []ast.Statement) object.Object {
	var res object.Object

	for _, s := range ss {
		res = Eval(s)

		// if res is a Return or an Error, stop evaluating and return it immediately
		if r, ok := res.(*object.Return); ok {
			return r
		}
		if e, ok := res.(*object.Error); ok {
			return e
		}
	}

	return res
}
