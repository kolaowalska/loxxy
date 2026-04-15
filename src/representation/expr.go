package representation

import (
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

// Expr - Base interface for all expression nodes.
type Expr interface {
	exprNode() // dummy method
}

type Binary struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (b *Binary) exprNode() {}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) exprNode() {}

type Literal struct {
	Value any
}

func (l *Literal) exprNode() {}

type Unary struct {
	Operator scanner.Token
	Right    Expr
}

func (u *Unary) exprNode() {}

type Variable struct {
	Name scanner.Token
}

func (v *Variable) exprNode() {}

type Assign struct {
	Name  scanner.Token
	Value Expr
}

func (a *Assign) exprNode() {}
