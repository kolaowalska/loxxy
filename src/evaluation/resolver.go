package evaluation

import (
	"fmt"

	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type Resolver struct {
	interpreter *Interpreter
	scopes      []map[string]bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
		scopes:      make([]map[string]bool, 0),
	}
}

// --- Scope Management ---
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

		// PERSON 2
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

		// PERSON 2
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
