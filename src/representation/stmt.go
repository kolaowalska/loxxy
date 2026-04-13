package representation

type Stmt interface {
	stmtNode()
}

// TODO

type Expression struct{}
type Print struct{}
type Var struct{}
type Block struct{}
