package evaluation

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type Interpreter struct {
	environment *Environment
	Stdout      io.Writer
	LastValue   any
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(nil),
		Stdout:      os.Stdout,
	}
}

// TODO: should return RuntimeError

func (i *Interpreter) Interpret(statements []representation.Stmt) error {
	for _, statement := range statements {
		err := i.Execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) Execute(stmt representation.Stmt) error {
	switch s := stmt.(type) {
	case *representation.If:
		condition, err := i.Evaluate(s.Condition)
		if err != nil {
			return err
		}
		if isTruthy(condition) {
			err := i.Execute(s.ThenBranch)
			if err != nil {
				return err
			}
		} else if s.ElseBranch != nil {
			err := i.Execute(s.ElseBranch)
			if err != nil {
				return err
			}
		}
		return nil

	case *representation.Print:
		value, err := i.Evaluate(s.Expression)
		if err != nil {
			return err
		}
		i.LastValue = value
		_, _ = fmt.Fprintln(i.Stdout, stringify(value))
		return nil

	case *representation.Expression:
		value, err := i.Evaluate(s.Expression)
		i.LastValue = value
		return err

	case *representation.Var:
		var value any = nil
		var err error
		if s.Initializer != nil {
			value, err = i.Evaluate(s.Initializer)
			if err != nil {
				return err
			}
		}
		i.environment.Define(s.Name.Lexeme, value)
		return nil

	case *representation.Block:
		return i.executeBlock(s.Statements, NewEnvironment(i.environment))

	case *representation.While:
		for {
			cond, err := i.Evaluate(s.Expression)
			if err != nil {
				return err
			}
			if isTruthy(cond) {
				i.Execute(s.Statement)
			} else {
				return nil
			}
		}
	}
	return fmt.Errorf("unknown statement type: %T", stmt)
}

func (i *Interpreter) executeBlock(statements []representation.Stmt, environment *Environment) error {
	previous := i.environment

	// try { ... } finally { ... }
	defer func() {
		i.environment = previous
	}()

	i.environment = environment

	for _, statement := range statements {
		err := i.Execute(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) Evaluate(expr representation.Expr) (any, error) {
	switch e := expr.(type) {

	case *representation.Literal:
		return e.Value, nil

	case *representation.Grouping:
		return i.Evaluate(e.Expression)

	case *representation.Variable:
		return i.environment.Get(e.Name)

	case *representation.Assign:
		value, err := i.Evaluate(e.Value)
		if err != nil {
			return nil, err
		}
		err = i.environment.Assign(e.Name, value)
		if err != nil {
			return nil, err
		}
		return value, nil

	case *representation.Logical:
		left, err := i.Evaluate(e.Left)
		if err != nil {
			return nil, err
		}

		if e.Operator.TokenType == scanner.OR {
			if isTruthy(left) {
				return left, nil
			}
		} else {
			if !isTruthy(left) {
				return left, nil
			}
		}

		return i.Evaluate(e.Right)

	case *representation.Unary:
		right, err := i.Evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.TokenType {

		case scanner.MINUS:
			err := checkNumberOperand(e.Operator, right)
			if err != nil {
				return nil, err
			}
			return -right.(float64), nil

		case scanner.BANG:
			return !isTruthy(right), nil

		default:
			return nil, fmt.Errorf("it's not supposed to go there, error in func Evaluate in Unary case - unknown operator")

		}

	case *representation.Binary:
		left, err := i.Evaluate(e.Left)
		if err != nil {
			return nil, err
		}

		right, err := i.Evaluate(e.Right)
		if err != nil {
			return nil, err
		}

		switch e.Operator.TokenType {

		case scanner.MINUS:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) - right.(float64), nil

		case scanner.SLASH:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			if right.(float64) == 0 { //TODO: test if needed
				return nil, newRuntimeError(e.Operator, "Cannot divide by zero.")
			}
			return left.(float64) / right.(float64), nil

		case scanner.STAR:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
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
			return nil, newRuntimeError(e.Operator, "Operands must be two numbers or two strings.")

		case scanner.GREATER:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) > right.(float64), nil

		case scanner.GREATER_EQUAL:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) >= right.(float64), nil

		case scanner.LESS:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) < right.(float64), nil

		case scanner.LESS_EQUAL:
			err := checkNumberOperands(e.Operator, left, right)
			if err != nil {
				return nil, err
			}
			return left.(float64) <= right.(float64), nil

		case scanner.BANG_EQUAL:
			return left != right, nil

		case scanner.EQUAL_EQUAL:
			return left == right, nil

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

func stringify(val any) string {
	if val == nil {
		return "nil"
	}

	if num, ok := val.(float64); ok {
		text := fmt.Sprintf("%v", num)
		if strings.HasSuffix(text, ".0") {
			return text[0 : len(text)-1] // -2?
		}
		return text
	}
	return fmt.Sprintf("%v", val)
}
