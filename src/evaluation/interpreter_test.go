package evaluation_test

import (
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name          string
		expr          representation.Expr
		expected      any
		expectedError bool
	}{
		// Literals
		{"Literal number", &representation.Literal{Value: 42.0}, 42.0, false},
		{"Literal string", &representation.Literal{Value: "hello"}, "hello", false},
		{"Literal boolean", &representation.Literal{Value: true}, true, false},
		{"Literal Nil", &representation.Literal{Value: nil}, nil, false},

		// Grouping
		{"Grouping", &representation.Grouping{Expression: &representation.Literal{Value: 10.0}}, 10.0, false},

		// Unary
		{"Unary bang true", &representation.Unary{Operator: scanner.Token{TokenType: scanner.BANG}, Right: &representation.Literal{Value: true}}, false, false},
		{"Unary bang false", &representation.Unary{Operator: scanner.Token{TokenType: scanner.BANG}, Right: &representation.Literal{Value: false}}, true, false},
		{"Unary bang nil", &representation.Unary{Operator: scanner.Token{TokenType: scanner.BANG}, Right: &representation.Literal{Value: nil}}, true, false},
		{"Unary bang number", &representation.Unary{Operator: scanner.Token{TokenType: scanner.BANG}, Right: &representation.Literal{Value: 5.0}}, false, false},
		{"Unary minus number", &representation.Unary{Operator: scanner.Token{TokenType: scanner.MINUS}, Right: &representation.Literal{Value: 5.0}}, -5.0, false},
		{"Unary minus string", &representation.Unary{Operator: scanner.Token{TokenType: scanner.MINUS}, Right: &representation.Literal{Value: "str"}}, nil, true},
		{"Unary bang zero", &representation.Unary{Operator: scanner.Token{TokenType: scanner.BANG}, Right: &representation.Literal{Value: 0.0}}, false, false},
		{"Unary bang empty string", &representation.Unary{Operator: scanner.Token{TokenType: scanner.BANG}, Right: &representation.Literal{Value: ""}}, false, false},

		// Binary - arithmetics
		{"Binary add numbers", &representation.Binary{Left: &representation.Literal{Value: 2.0}, Operator: scanner.Token{TokenType: scanner.PLUS}, Right: &representation.Literal{Value: 3.0}}, 5.0, false},
		{"Binary add strings", &representation.Binary{Left: &representation.Literal{Value: "ala"}, Operator: scanner.Token{TokenType: scanner.PLUS}, Right: &representation.Literal{Value: "kot"}}, "alakot", false},
		{"Binary add mixed", &representation.Binary{Left: &representation.Literal{Value: 2.0}, Operator: scanner.Token{TokenType: scanner.PLUS}, Right: &representation.Literal{Value: "kot"}}, nil, true},
		{"Binary subtract", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.MINUS}, Right: &representation.Literal{Value: 3.0}}, 2.0, false},
		{"Binary multiply", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.STAR}, Right: &representation.Literal{Value: 3.0}}, 15.0, false},
		{"Binary divide", &representation.Binary{Left: &representation.Literal{Value: 6.0}, Operator: scanner.Token{TokenType: scanner.SLASH}, Right: &representation.Literal{Value: 2.0}}, 3.0, false},
		{"Binary divide by zero", &representation.Binary{Left: &representation.Literal{Value: 6.0}, Operator: scanner.Token{TokenType: scanner.SLASH}, Right: &representation.Literal{Value: 0.0}}, nil, true},

		// Binary - logic operators
		{"Binary greater", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.GREATER}, Right: &representation.Literal{Value: 3.0}}, true, false},
		{"Binary greater equal", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.GREATER_EQUAL}, Right: &representation.Literal{Value: 5.0}}, true, false},
		{"Binary less", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.LESS}, Right: &representation.Literal{Value: 8.0}}, true, false},
		{"Binary less equal", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.LESS_EQUAL}, Right: &representation.Literal{Value: 5.0}}, true, false},
		{"Binary less equal mixed", &representation.Binary{Left: &representation.Literal{Value: "ala"}, Operator: scanner.Token{TokenType: scanner.LESS_EQUAL}, Right: &representation.Literal{Value: 5.0}}, nil, true},
		{"Binary equal equal numbers", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL}, Right: &representation.Literal{Value: 5.0}}, true, false},
		{"Binary equal equal strings", &representation.Binary{Left: &representation.Literal{Value: "ala"}, Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL}, Right: &representation.Literal{Value: "kot"}}, false, false},
		{"Binary equal equal booleans", &representation.Binary{Left: &representation.Literal{Value: false}, Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL}, Right: &representation.Literal{Value: false}}, true, false},
		{"Binary equal equal mixed types", &representation.Binary{Left: &representation.Literal{Value: "abc"}, Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL}, Right: &representation.Literal{Value: 10.0}}, false, false},
		{"Binary equal equal nil and number", &representation.Binary{Left: &representation.Literal{Value: nil}, Operator: scanner.Token{TokenType: scanner.EQUAL_EQUAL}, Right: &representation.Literal{Value: 0.0}}, false, false},
		{"Binary bang equal numbers", &representation.Binary{Left: &representation.Literal{Value: 5.0}, Operator: scanner.Token{TokenType: scanner.BANG_EQUAL}, Right: &representation.Literal{Value: 3.0}}, true, false},
		{"Binary bang equal strings", &representation.Binary{Left: &representation.Literal{Value: "ala"}, Operator: scanner.Token{TokenType: scanner.BANG_EQUAL}, Right: &representation.Literal{Value: "kot"}}, true, false},
		{"Binary bang equal booleans", &representation.Binary{Left: &representation.Literal{Value: false}, Operator: scanner.Token{TokenType: scanner.BANG_EQUAL}, Right: &representation.Literal{Value: false}}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluation.Evaluate(tt.expr)

			if tt.expectedError {
				if err == nil {
					t.Errorf("evaluation.Evaluate(%v) expected error but got nil", tt.expr)
				}
			} else {
				if err != nil {
					t.Errorf("evaluation.Evaluate(%v) returned an error and it shouldn't %v", tt.expr, err)
				}
				if result != tt.expected {
					t.Errorf("evaluation.Evaluate(%v) expected %v but got %v", tt.expr, tt.expected, result)
				}
			}
		})
	}

}
