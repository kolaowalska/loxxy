package resolving

import (
	"github.com/kolaowalska/loxxy/src/evaluation"
	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
)

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
	FunctionTypeInitializer
	FunctionTypeMethod
)

type ErrorReporter interface {
	Error(line int, message string)
	TokenError(t scanner.Token, message string)
}

type Resolver struct {
	interpreter     *evaluation.Interpreter
	scopes          []map[string]bool
	currentClass    ClassType
	currentFunction FunctionType
	reporter        ErrorReporter
}

func NewResolver(interpreter *evaluation.Interpreter, reporter ErrorReporter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          make([]map[string]bool, 0),
		currentClass:    ClassTypeNone,
		currentFunction: FunctionTypeNone,
		reporter:        reporter,
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		r.reporter.TokenError(name, "already a variable with this name in this scope")
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name scanner.Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes[len(r.scopes)-1][name.Lexeme] = true // true means "ready to use"
}

func (r *Resolver) ResolveStatements(statements []representation.Stmt) error {
	for _, stmt := range statements {
		err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt representation.Stmt) error {
	switch s := stmt.(type) {
	case *representation.Block:
		r.beginScope()
		err := r.ResolveStatements(s.Statements)
		r.endScope()
		return err

	case *representation.Var:
		r.declare(s.Name)
		if s.Initializer != nil {
			err := r.resolveExpr(s.Initializer)
			if err != nil {
				return err
			}
		}
		r.define(s.Name)
		return nil

	case *representation.Function:
		r.declare(s.Name)
		r.define(s.Name)
		return r.resolveFunction(s, FunctionTypeFunction)

	case *representation.Expression:
		return r.resolveExpr(s.Expression)

	case *representation.If:
		err := r.resolveExpr(s.Condition)
		if err != nil {
			return err
		}
		err = r.resolveStmt(s.ThenBranch)

		if err != nil {
			return err
		}
		if s.ElseBranch != nil {
			return r.resolveStmt(s.ElseBranch)
		}
		return nil

	case *representation.Print:
		return r.resolveExpr(s.Expression)

	case *representation.Return:
		if r.currentFunction == FunctionTypeNone {
			r.reporter.TokenError(s.Keyword, "can't return from top-level code")
		}
		if s.Value != nil {
			if r.currentFunction == FunctionTypeInitializer {
				r.reporter.TokenError(s.Keyword, "can't return a value from an initializer")
			}
			return r.resolveExpr(s.Value)
		}

	case *representation.While:
		err := r.resolveExpr(s.Condition)
		if err != nil {
			return err
		}
		return r.resolveStmt(s.Body)

	case *representation.Class:
		r.declare(s.Name)
		r.define(s.Name)

		enclosingClass := r.currentClass
		r.currentClass = ClassTypeClass

		r.beginScope()
		r.scopes[len(r.scopes)-1]["this"] = true

		for _, method := range s.Methods {
			declaration := FunctionTypeMethod
			if method.Name.Lexeme == "init" {
				declaration = FunctionTypeInitializer
			}
			err := r.resolveFunction(method, declaration)
			if err != nil {
				return err
			}
		}

		r.endScope()
		r.currentClass = enclosingClass

		return nil
	}
	return nil
}

func (r *Resolver) resolveExpr(expr representation.Expr) error {
	switch e := expr.(type) {
	case *representation.Variable:
		if len(r.scopes) != 0 {
			scope := r.scopes[len(r.scopes)-1]
			if ready, exists := scope[e.Name.Lexeme]; exists && !ready {
				r.reporter.TokenError(e.Name, "can't read local variable in its own initializer.")
			}
		}
		r.resolveLocal(e, e.Name)

	case *representation.Assign:
		err := r.resolveExpr(e.Value)
		if err != nil {
			return err
		}
		r.resolveLocal(e, e.Name)
		return nil

	case *representation.Binary:
		err := r.resolveExpr(e.Left)
		if err != nil {
			return err
		}
		return r.resolveExpr(e.Right)

	case *representation.Call:
		err := r.resolveExpr(e.Callee)
		if err != nil {
			return err
		}
		for _, arg := range e.Args {
			err := r.resolveExpr(arg)
			if err != nil {
				return err
			}
		}
		return nil

	case *representation.Grouping:
		return r.resolveExpr(e.Expression)

	case *representation.Literal:
		return nil

	case *representation.Logical:
		err := r.resolveExpr(e.Left)
		if err != nil {
			return err
		}
		return r.resolveExpr(e.Right)

	case *representation.Unary:
		return r.resolveExpr(e.Right)

	case *representation.Get:
		return r.resolveExpr(e.Object)

	case *representation.Set:
		err := r.resolveExpr(e.Value)
		if err != nil {
			return err
		}
		return r.resolveExpr(e.Object)

	case *representation.This:
		if r.currentClass == ClassTypeNone {
			r.reporter.TokenError(e.Keyword, "can't use 'this' outside of a class.")
		}
		r.resolveLocal(e, e.Keyword)
	}
	return nil
}

func (r *Resolver) resolveLocal(expr representation.Expr, name scanner.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, exists := r.scopes[i][name.Lexeme]; exists {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveFunction(function *representation.Function, fType FunctionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = fType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	err := r.ResolveStatements(function.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
	return err
}
