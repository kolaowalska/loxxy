package resolving

import (
	"fmt"

	"github.com/kolaowalska/loxxy/src/evaluation"
	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type ClassType int

const (
	ClassTypeNone ClassType = iota
	ClassTypeClass
)

type Resolver struct {
	interpreter  *evaluation.Interpreter
	scopes       []map[string]bool
	currentClass ClassType
}

func NewResolver(interpreter *evaluation.Interpreter) *Resolver {
	return &Resolver{
		interpreter:  interpreter,
		scopes:       make([]map[string]bool, 0),
		currentClass: ClassTypeNone,
	}
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name scanner.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}
	scope := r.scopes[len(r.scopes)-1]
	if _, exists := scope[name.Lexeme]; exists {
		return fmt.Errorf("already a variable with this name in this scope")
	}
	scope[name.Lexeme] = false
	return nil
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
		err := r.declare(s.Name)
		if err != nil {
			return err
		}
		if s.Initializer != nil {
			err = r.resolveExpr(s.Initializer)
			if err != nil {
				return err
			}
		}
		r.define(s.Name)
		return nil

	case *representation.Function:
		err := r.declare(s.Name)
		if err != nil {
			return err
		}
		r.define(s.Name)
		return r.resolveFunction(s)

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
		if s.Value != nil {
			return r.resolveExpr(s.Value)
		}
		return nil

	case *representation.While:
		err := r.resolveExpr(s.Condition)
		if err != nil {
			return err
		}
		return r.resolveStmt(s.Body)

	case *representation.Class:
		err := r.declare(s.Name)
		if err != nil {
			return err
		}
		r.define(s.Name)

		enclosingClass := r.currentClass
		r.currentClass = ClassTypeClass

		r.beginScope()
		r.scopes[len(r.scopes)-1]["this"] = true

		for _, method := range s.Methods {
			err := r.resolveFunction(method)
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
				return fmt.Errorf("can't read local variable in its own initializer")
			}
		}
		r.resolveLocal(e, e.Name)
		return nil

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
			return fmt.Errorf("can't use 'this' outside of a class")
		}
		r.resolveLocal(e, e.Keyword)
		return nil
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

func (r *Resolver) resolveFunction(function *representation.Function) error {
	r.beginScope()
	for _, param := range function.Params {
		err := r.declare(param)
		if err != nil {
			return err
		}
		r.define(param)
	}
	err := r.ResolveStatements(function.Body)
	r.endScope()
	return err
}
