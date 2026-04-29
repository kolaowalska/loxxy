package evaluation

import scanner "github.com/kolaowalska/loxxy/src/scanning"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{enclosing: enclosing, values: make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}
func (e *Environment) Get(name scanner.Token) (any, error) {
	if val, ok := e.values[name.Lexeme]; ok {
		return val, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, newRuntimeError(name, "undefined variable '"+name.Lexeme+"'.")
}
func (e *Environment) Assign(name scanner.Token, value any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return newRuntimeError(name, "undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.enclosing
	}
	return environment
}

func (e *Environment) GetAt(distance int, name string) (any, error) {
	return e.ancestor(distance).values[name], nil
}

func (e *Environment) AssignAt(distance int, name scanner.Token, value any) error {
	e.ancestor(distance).values[name.Lexeme] = value
	return nil
}
