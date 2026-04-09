package evaluation

import scanner "github.com/kolaowalska/loxxy/src/scanning"

type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func (r *RuntimeError) Error() string {
	return r.Message
}

func NewRuntimeError(operator scanner.Token, message string) *RuntimeError {
	return &RuntimeError{Token: operator, Message: message}
}

func checkNumberOperand(operator scanner.Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return NewRuntimeError(operator, "operand must be a number.")
}

func checkNumberOperands(operator scanner.Token, left any, right any) error {
	_, leftIsNumber := left.(float64)
	_, rightIsNumber := right.(float64)
	if leftIsNumber && rightIsNumber {
		return nil
	}
	return NewRuntimeError(operator, "operands must be numbers.")
}
