package representation

import scanner "github.com/kolaowalska/loxxy/src/scanning"

type Stmt interface {
	stmtNode()
}

// TODO: implement methods

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
