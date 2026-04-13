package representation

type Stmt interface {
	stmtNode()
}

// TODO: implement methods

type Expression struct{}
type Print struct{}
type Var struct{}
type Block struct{}
