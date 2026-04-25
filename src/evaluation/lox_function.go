package evaluation

import (
	"errors"

	"github.com/kolaowalska/loxxy/src/representation"
)

type LoxFunction struct {
	declaration *representation.Function
	closure     *Environment
}

func NewLoxFunction(declaration *representation.Function, closure *Environment) *LoxFunction {
	return &LoxFunction{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.Params)
}

func (f *LoxFunction) Call(i *Interpreter, arguments []any) (any, error) {
	environment := NewEnvironment(f.closure)

	for j, param := range f.declaration.Params {
		environment.Define(param.Lexeme, arguments[j])
	}

	err := i.executeBlock(f.declaration.Body, environment)
	if err != nil {
		// TODO: double-check cursed type error
		if ret, ok := errors.AsType[*ReturnValue](err); ok {
			return ret.Value, nil
		}
		return nil, err
	}
	return nil, nil
}

func (f *LoxFunction) String() string {
	return "<fn " + f.declaration.Name.Lexeme + ">"
}
