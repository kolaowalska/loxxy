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
	globals     *Environment
	locals      map[representation.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	globals.Define("clock", &NativeClock{})
	return &Interpreter{
		environment: globals,
		Stdout:      os.Stdout,
		globals:     globals,
		locals:      make(map[representation.Expr]int),
	}
}

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
		_, _ = fmt.Fprintln(i.Stdout, stringify(value))
		return nil

	case *representation.Expression:
		_, err := i.Evaluate(s.Expression)
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

	case *representation.Return:
		var value any = nil
		var err error
		if s.Value != nil {
			value, err = i.Evaluate(s.Value)
			if err != nil {
				return err
			}
		}
		return &ReturnValue{Value: value}

	case *representation.While:
		for {
			cond, err := i.Evaluate(s.Condition)
			if err != nil {
				return err
			}
			if isTruthy(cond) {
				err := i.Execute(s.Body)
				if err != nil {
					return err
				}
			} else {
				return nil
			}
		}
	case *representation.Function:
		function := NewLoxFunction(s, i.environment)
		i.environment.Define(s.Name.Lexeme, function)
		return nil
	}

	return fmt.Errorf("unknown statement type: %T", stmt)
}

func (i *Interpreter) Resolve(expr representation.Expr, depth int) {
	i.locals[expr] = depth
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
		return i.lookupVariable(e.Name, e)

	case *representation.Assign:
		value, err := i.Evaluate(e.Value)
		if err != nil {
			return nil, err
		}

		distance, ok := i.locals[e]
		if ok {
			err = i.environment.AssignAt(distance, e.Name, value)
		} else {
			err = i.globals.Assign(e.Name, value)
		}
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

	case *representation.Call:
		callee, err := i.Evaluate(e.Callee)
		if err != nil {
			return nil, err
		}

		var arguments []any
		for _, argExpr := range e.Args {
			arg, err := i.Evaluate(argExpr)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
		}

		function, ok := callee.(LoxCallable)
		if !ok {
			return nil, newRuntimeError(e.Paren, "can only call functions and classes.")
		}
		if len(arguments) != function.Arity() {
			return nil, newRuntimeError(e.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)))
		}

		return function.Call(i, arguments)
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

func (i *Interpreter) lookupVariable(name scanner.Token, expr representation.Expr) (any, error) {
	distance, ok := i.locals[expr]
	if ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.globals.Get(name)
	}
}
