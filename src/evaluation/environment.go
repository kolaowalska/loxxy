package evaluation

import (
	"github.com/kolaowalska/loxxy/src/scanning"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing: enclosing, values: make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	// TODO
}
func (e *Environment) Get(name scanner.Token) (any, error) {
	// TODO
	return nil, newRuntimeError(name, "undefined variable '"+name.Lexeme+"'.")
}
func (e *Environment) Assign(name scanner.Token, value any) error {
	// TODO
	return newRuntimeError(name, "undefined variable '"+name.Lexeme+"'.")
}
