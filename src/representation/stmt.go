package representation

import scanner "github.com/kolaowalska/loxxy/src/scanning"

type Stmt interface {
	stmtNode()
}

type Expression struct {
	Expression Expr
}

func (e *Expression) stmtNode() {}

type Print struct {
	Expression Expr
}

func (p *Print) stmtNode() {}

type Var struct {
	Name        scanner.Token
	Initializer Expr
}

func (v *Var) stmtNode() {}

type Block struct {
	Statements []Stmt
}

func (b *Block) stmtNode() {}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) stmtNode() {}

type While struct {
	Condition Expr
	Body      Stmt
}

func (w *While) stmtNode() {}

type Function struct {
	Name   scanner.Token
	Params []scanner.Token
	Body   []Stmt
}

func (f *Function) stmtNode() {}

type Return struct {
	Keyword scanner.Token
	Value   Expr
}

func (r *Return) stmtNode() {}

type Class struct {
	Name       scanner.Token
	Superclass *Variable
	Methods    []*Function
}

func (c *Class) stmtNode() {}

type Super struct {
	Keyword scanner.Token
	Method  scanner.Token
}

func (s *Super) exprNode() {}
