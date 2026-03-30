package representation

import (
	"testing"

	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func TestAstPrinter(t *testing.T) {
	expression := &Binary{
		Left: &Unary{
			Operator: scanner.Token{scanner.MINUS, "-", nil, 1},
			Right:    &Literal{123},
		},
		Operator: scanner.Token{scanner.STAR, "*", nil, 1},
		Right: &Grouping{
			Expression: &Literal{45.67},
		},
	}
	result := Print(expression)
	expected := "(* (- 123) (group 45.67))"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
