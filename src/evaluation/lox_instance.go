package evaluation

import scanner "github.com/kolaowalska/loxxy/src/scanning"

type LoxInstance struct {
	Class  *LoxClass
	Fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		Class:  class,
		Fields: make(map[string]any),
	}
}

func (i *LoxInstance) String() string {
	return i.Class.Name + " instance"
}

func (i *LoxInstance) Get(name scanner.Token) (any, error) {
	if val, ok := i.Fields[name.Lexeme]; ok {
		return val, nil
	}
	if method := i.Class.FindMethod(name.Lexeme); method != nil {
		return method.Bind(i), nil
	}
	return nil, newRuntimeError(name, "undefined property '"+name.Lexeme+"'.")
}

func (i *LoxInstance) Set(name scanner.Token, value any) {
	i.Fields[name.Lexeme] = value
}
