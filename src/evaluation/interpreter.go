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
	Hook        DebugHook   // nil when not debugging
	callStack   []CallFrame // maintained during execution
}

func NewInterpreter() *Interpreter {
	builtins := NewEnvironment(nil)
	builtins.Define("clock", &NativeClock{})
	globals := NewEnvironment(builtins)
	globals.isGlobalScope = true
	return &Interpreter{
		environment: globals,
		Stdout:      os.Stdout,
		globals:     globals,
		locals:      make(map[representation.Expr]int),
	}
}

func (i *Interpreter) Interpret(statements []representation.Stmt) (err error) {
	i.callStack = append(i.callStack, CallFrame{Name: "<script>", Line: 0})
	defer func() { i.callStack = i.callStack[:len(i.callStack)-1] }()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("internal interpreter panic: %v", r)
		}
	}()

	for _, statement := range statements {
		executeError := i.Execute(statement)
		if executeError != nil {
			return executeError
		}
	}
	return nil
}

func (i *Interpreter) Execute(stmt representation.Stmt) error {
	if i.Hook != nil {
		if line := stmtLine(stmt); line > 0 {
			if len(i.callStack) > 0 {
				i.callStack[len(i.callStack)-1].Line = line
			}
			i.Hook.OnStatement(line, append([]CallFrame{}, i.callStack...), i.environment)
		}
	}

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
		function := NewLoxFunction(s, i.environment, false)
		i.environment.Define(s.Name.Lexeme, function)
		return nil

	case *representation.Class:
		var superclass any = nil
		if s.Superclass != nil {
			var err error
			superclass, err = i.Evaluate(s.Superclass)
			if err != nil {
				return err
			}

			if _, ok := superclass.(*LoxClass); !ok {
				return newRuntimeError(s.Superclass.Name, "superclass must be a class.")
			}
		}

		i.environment.Define(s.Name.Lexeme, nil)

		if s.Superclass != nil {
			i.environment = NewEnvironment(i.environment)
			i.environment.Define("super", superclass)
		}

		methods := make(map[string]*LoxFunction)
		for _, method := range s.Methods {
			isInit := method.Name.Lexeme == "init"
			function := NewLoxFunction(method, i.environment, isInit)
			methods[method.Name.Lexeme] = function
		}

		var superclassPtr *LoxClass = nil
		if superclass != nil {
			superclassPtr = superclass.(*LoxClass)
		}

		class := &LoxClass{Name: s.Name.Lexeme, Superclass: superclassPtr, Methods: methods}

		if superclassPtr != nil {
			i.environment = i.environment.enclosing
		}

		err := i.environment.Assign(s.Name, class)
		if err != nil {
			return err
		}
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
			return nil, fmt.Errorf("error in func Evaluate in Unary case - unknown operator")

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
				return nil, newRuntimeError(e.Operator, "cannot divide by zero.")
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
			return nil, newRuntimeError(e.Operator, "operands must be two numbers or two strings.")

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

	case *representation.This:
		return i.lookupVariable(e.Keyword, e)

	case *representation.Get:
		object, err := i.Evaluate(e.Object)
		if err != nil {
			return nil, err
		}
		if instance, ok := object.(*LoxInstance); ok {
			return instance.Get(e.Name)
		}
		return nil, newRuntimeError(e.Name, "only instances have properties.")

	case *representation.Set:
		object, err := i.Evaluate(e.Object)
		if err != nil {
			return nil, err
		}
		if instance, ok := object.(*LoxInstance); ok {
			value, err := i.Evaluate(e.Value)
			if err != nil {
				return nil, err
			}
			instance.Set(e.Name, value)
			return value, nil
		}
		return nil, newRuntimeError(e.Name, "only instances have fields.")
	case *representation.Super:
		distance := i.locals[e]

		superclassVal, _ := i.environment.GetAt(distance, "super")
		superclass := superclassVal.(*LoxClass)

		objectVal, _ := i.environment.GetAt(distance-1, "this")
		object := objectVal.(*LoxInstance)

		method := superclass.FindMethod(e.Method.Lexeme)
		if method == nil {
			return nil, newRuntimeError(e.Method, "undefined property '"+e.Method.Lexeme+"'.")
		}

		return method.Bind(object), nil
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
	}

	return i.globals.Get(name)
}

func (i *Interpreter) ClearLocals() {
	i.locals = make(map[representation.Expr]int)
}

// function to extract the best available line from any statement
func stmtLine(stmt representation.Stmt) int {
	switch s := stmt.(type) {
	case *representation.Var:
		return s.Name.Line
	case *representation.Function:
		return s.Name.Line
	case *representation.If:
		return exprLine(s.Condition)
	case *representation.While:
		return exprLine(s.Condition)
	case *representation.Print:
		return exprLine(s.Expression)
	case *representation.Return:
		return s.Keyword.Line
	case *representation.Expression:
		return exprLine(s.Expression)
	case *representation.Class:
		return s.Name.Line
	}
	return 0
}

func exprLine(expr representation.Expr) int {
	switch e := expr.(type) {
	case *representation.Assign:
		return e.Name.Line
	case *representation.Binary:
		return e.Operator.Line
	case *representation.Call:
		return e.Paren.Line
	case *representation.Get:
		return e.Name.Line
	case *representation.Logical:
		return e.Operator.Line
	case *representation.Set:
		return e.Name.Line
	case *representation.Super:
		return e.Keyword.Line
	case *representation.This:
		return e.Keyword.Line
	case *representation.Unary:
		return e.Operator.Line
	case *representation.Variable:
		return e.Name.Line
	}
	return 0
}
