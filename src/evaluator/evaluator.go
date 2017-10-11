package evaluator

import (
	"ast"
	"fmt"
	"object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		// special case for Program, need to unwrap Return
		res := evalStatements(node.Statements, env)
		if res.Type() == object.TYPE_RETURN {
			return res.(*object.Return).Value
		} else {
			return res
		}
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.IntValue}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: node.BoolValue}
	case *ast.PrefixExpression:
		right := Eval(node.Expression, env)
		return evalPrefix(node.Operator, right, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfix(node.Operator, left, right, env)
	case *ast.IfExpression:
		condition := Eval(node.Condition, env)

		pred, err := convertToBool(condition)
		if err != nil {
			return err
		}

		if pred {
			return Eval(node.Consequence, env)
		} else {
			if node.Alternative == nil {
				return &object.Null{}
			} else {
				return Eval(node.Alternative, env)
			}
		}
	case *ast.BlockStatement:
		return evalStatements(node.Statements, env)
	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.Value, env)}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if value.Type() == object.TYPE_ERROR {
			return value
		}
		env.Set(node.Ident.Name, value)
		return value
	case *ast.Identifier:
		if value, ok := env.Get(node.Name); !ok {
			return newError("unknown identifier: %s", node.Name)
		} else {
			return value
		}
	case *ast.FunctionExpression:
		params := []string{}
		for _, p := range node.Params {
			params = append(params, p.Name)
		}
		return &object.Function{Params: params, Body: node.Body, Env: env}
	case *ast.CallExpression:
		c := Eval(node.Function, env)
		if c.Type() == object.TYPE_ERROR {
			return c
		}
		if c.Type() != object.TYPE_FUNCTION {
			return newError("non callable object is used: %s", c.Inspect())
		}

		f := c.(*object.Function)
		e2 := f.Env.NewLinkedEnvironment()
		for i, arg := range node.Arguments {
			actual := Eval(arg, env)
			if actual.Type() == object.TYPE_ERROR {
				return actual
			}

			e2.Set(f.Params[i], actual)
		}

		value := Eval(f.Body, e2)
		if value.Type() == object.TYPE_RETURN {
			return value.(*object.Return).Value
		} else {
			return value
		}
	default:
		return newError("unhandled case %T", node)
	}
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func evalPrefix(operator string, operand object.Object, env *object.Environment) object.Object {
	if operand.Type() == object.TYPE_ERROR {
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

func evalInfix(operator string, left object.Object, right object.Object, env *object.Environment) object.Object {
	if left.Type() == object.TYPE_ERROR {
		return left
	}
	if right.Type() == object.TYPE_ERROR {
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

func evalStatements(ss []ast.Statement, env *object.Environment) object.Object {
	var res object.Object

	for _, s := range ss {
		res = Eval(s, env)

		// if res is a Return or an Error, stop evaluating and return it immediately
		if res.Type() == object.TYPE_RETURN || res.Type() == object.TYPE_ERROR {
			return res
		}
	}

	return res
}
