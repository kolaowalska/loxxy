package parsing

import (
	"testing"

	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type MockReporter struct {
	HadError    bool
	LastMessage string
}

func (m *MockReporter) error(t scanner.Token, message string) {
	m.HadError = true
	m.LastMessage = message
}
func (m *MockReporter) Error(line int, message string) {
	m.HadError = true
	m.LastMessage = message
}

//TODO: revise tests after full implementation

func TestParser_ValidExpressions(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		// Literals
		{"number literal", "123", "123"},
		{"string literal", "\"hello\"", "hello"},
		{"string literal with spaces", "\"hello world\"", "hello world"},
		{"empty string", "\"\"", ""},
		{"true literal", "true", "true"},
		{"false literal", "false", "false"},
		{"nil literal", "nil", "nil"},

		// Unary
		{"unary minus", "-123", "(- 123)"},
		{"unary bang", "!true", "(! true)"},
		{"nested unary", "!!false", "(! (! false))"},

		// Comparison
		{"greater comparison", "1 > 2", "(> 1 2)"},
		{"less comparison", "1 < 2", "(< 1 2)"},
		{"greater or equal comparison", "1 >= 2", "(>= 1 2)"},
		{"less or equal comparison", "1 <= 2", "(<= 1 2)"},

		// Associativity, Term (- +), Factor (/ *)
		{"left associativity same precedence - subtraction, addition", "1 + 2 - 1", "(- (+ 1 2) 1)"},
		{"left associativity same precedence - multiplication, division", "10 / 2 * 5", "(* (/ 10 2) 5)"},
		{"left associativity with multiple same precedence operations", "1 + 2 + 3 + 4 + 5", "(+ (+ (+ (+ 1 2) 3) 4) 5)"},

		// Precedence
		{"unary over multiplication", "-3 * 5", "(* (- 3) 5)"},
		{"multiplication over addition", "1 + 2 * 3", "(+ 1 (* 2 3))"},
		{"addition over comparison", "1 + 2 < 4", "(< (+ 1 2) 4)"},
		{"comparison over equality", "1 < 2 == 3 > 4", "(== (< 1 2) (> 3 4))"},

		// Grouping
		{"grouping overrides precedence", "(1 + 2) * 3", "(* (group (+ 1 2)) 3)"},
		{"nested grouping", "((12 + 31) - 4)", "(group (- (group (+ 12 31)) 4))"},

		// Misc
		{"nil equality", "nil == false", "(== nil false)"},
		{"string and number combination", "\"123\" + 456", "(+ 123 456)"}, //I don't know if this does make sense, but it should go through at parsing stage
		{"mega test", "1 == 2 < 3 + 4 * 5", "(== 1 (< 2 (+ 3 (* 4 5))))"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockReporter{}

			s := scanner.NewScanner(tt.source, mock)
			tokens := s.ScanTokens()

			if mock.HadError {
				t.Fatalf("scanner reported an unexpected error for source: %q", tt.source)
			}

			p := NewParser(tokens, mock)
			expr := p.Parse()

			if mock.HadError {
				t.Fatalf("parser reported an unexpected error: %s", mock.LastMessage)
			}
			if expr == nil {
				t.Fatalf("parser returned nil expression: %q", tt.source)
			}

			result := representation.Print(expr)
			if result != tt.expected {
				t.Errorf("expected AST %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestParser_SyntaxErrors(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedError string
	}{
		{"missing expression after operator", "1 + ", "Expect expression."},
		{"starting with binary operator", "* 5", "Expect expression."},
		{"missing right parenthesis", "(1 + 2", "Expect ')' after expression."},
		{"only unary operator", "-", "Expect expression."},
		{"empty parenthesis", "()", "Expect expression."},
		{"error in nested grouping", "(1 + (2 * ))", "Expect expression."},
		{"missing left operand for equality", "== true", "Expect expression."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockReporter{}

			s := scanner.NewScanner(tt.source, mock)
			tokens := s.ScanTokens()

			p := NewParser(tokens, mock)
			p.Parse()

			if !mock.HadError {
				t.Errorf("expected parser to report an error for %q, but it didn't ", tt.source)
			}

			// If it bombs here, check error messages, bcs they do not exist when i'm writing this, so might be differences
			if mock.HadError && mock.LastMessage != tt.expectedError {
				t.Errorf("expected error message %q, got %q", tt.expectedError, mock.LastMessage)
			}
		})
	}
}
