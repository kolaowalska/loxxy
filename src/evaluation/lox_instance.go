package evaluation

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
