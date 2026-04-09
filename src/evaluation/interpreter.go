package evaluation

import (
	"fmt"

	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func Evaluate(expr representation.Expr) (any, error) {
	switch e := expr.(type) {

	case *representation.Literal:
		return e.Value, nil

	case *representation.Grouping:
		return Evaluate(e.Expression)

	case *representation.Unary:
		right, err := Evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.TokenType {

		case scanner.MINUS:
			return -right.(float64), nil

		case scanner.BANG:
			return !isTruthy(right), nil

		default:
			return nil, fmt.Errorf("it's not supposed to go there, error in func Evaluate in Unary case - unknown operator")

		}

	case *representation.Binary:
		left, err := Evaluate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := Evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.TokenType {

		case scanner.MINUS:
			return left.(float64) - right.(float64), nil

		case scanner.SLASH:
			if right.(float64) == 0 { //TODO: test if needed
				return nil, nil //TODO: insert runtime error on e.Operator, "Cannot divide by zero."
			}
			return left.(float64) / right.(float64), nil

		case scanner.STAR:
			return left.(float64) * right.(float64), nil

		case scanner.PLUS:
			if l, ok := left.(float64); ok {
				if r, ok := right.(float64); ok {
					return l + r, nil
				}
			}
			if l, ok := left.(string); ok {
				if r, ok := right.(string); ok {
					return l + r, nil
				}
			}
			return nil, nil //TODO: runtime error on e.Operator, "Operands must be two numbers or two strings."

		case scanner.GREATER:
			return left.(float64) > right.(float64), nil

		case scanner.GREATER_EQUAL:
			return left.(float64) >= right.(float64), nil

		case scanner.LESS:
			return left.(float64) < right.(float64), nil

		case scanner.LESS_EQUAL:
			return left.(float64) <= right.(float64), nil

		case scanner.BANG_EQUAL:
			return left.(float64) != right.(float64), nil

		case scanner.EQUAL_EQUAL:
			return left.(float64) == right.(float64), nil

		default:
			return nil, fmt.Errorf("it's not supposed to go there, error in func Evaluate in Binary case - unknown operator")
		}
	}

	return nil, fmt.Errorf("unknown expression type")
}

func isTruthy(obj any) bool {
	if obj == nil {
		return false
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	return true
}

// unnecessary due to go's two-step automatic check
//
//func isEqual(a any, b any) bool {
//	if a == nil && b == nil {
//		return true
//	}
//	if a == nil {
//		return false
//	}
//	return a == b
//}
