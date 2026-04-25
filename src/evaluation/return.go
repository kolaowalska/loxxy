package evaluation

type ReturnValue struct {
	Value any
}

func (r *ReturnValue) Error() string {
	return "return"
}
