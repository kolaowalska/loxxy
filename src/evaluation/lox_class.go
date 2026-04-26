package evaluation

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]*LoxFunction
}

func (c *LoxClass) String() string {
	return c.Name
}

func (c *LoxClass) Arity() int {
	if init := c.FindMethod("init"); init != nil {
		return init.Arity()
	}
	return 0
}

func (c *LoxClass) Call(i *Interpreter, arguments []any) (any, error) {
	instance := NewLoxInstance(c)
	if init := c.FindMethod("init"); init != nil {
		_, err := init.Bind(instance).Call(i, arguments)
		if err != nil {
			return nil, err
		}
	}
	return instance, nil
}

func (c *LoxClass) FindMethod(name string) *LoxFunction {
	if method, ok := c.Methods[name]; ok {
		return method
	}
	if c.Superclass != nil {
		return c.Superclass.Superclass.FindMethod(name)
	}
	return nil
}
