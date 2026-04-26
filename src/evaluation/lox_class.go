package evaluation

type LoxClass struct {
	Name    string
	Methods map[string]*LoxFunction
}

func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) Call(i *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(c)
	return instance, nil
}
