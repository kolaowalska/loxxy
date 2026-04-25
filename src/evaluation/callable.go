package evaluation

import "time"

type LoxCallable interface {
	Arity() int
	Call(interpreter *Interpreter, arguments []any) (any, error)
}

type NativeClock struct{}

func (n *NativeClock) Arity() int { return 0 }
func (n *NativeClock) Call(interpreter *Interpreter, arguments []any) (any, error) {
	return float64(time.Now().UnixNano()) / 1e9, nil
}
func (n *NativeClock) String() string { return "<native fn>" }
