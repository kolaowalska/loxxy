package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type TestReporter struct{}

func (r TestReporter) Error(line int, message string)                 {}
func (r TestReporter) TokenError(token scanner.Token, message string) {}

func TestPipeline(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{"Simple math", "1 + 2 * 3;", 7.0, false},
		{"Grouping", "(1 + 2) * 3;", 9.0, false},
		{"String concat", "\"ala\" + \"kot\";", "alakot", false},
		{"Comparison and booleans", "10 >= 5 == true;", true, false},
		{"Equality mixed", "\"ala\" == 123;", false, false},
		{"Complex test 1", "(-1 + 5) * 2 / 4;", 2.0, false},
		{"Complex test 2", "(5 - (3 - 1)) + -1;", 2.0, false},

		{"Subtract string from number", "\"ala\" - 123;", nil, true},
		{"Unary minus on string", "-\"ala\";", nil, true},
		{"Add booleans", "true + false;", nil, true},
		{"Divide by zero", "10 / 0;", nil, true},
		{"Greater comparison on strings", "\"a\" > \"b\";", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := TestReporter{}

			s := scanner.NewScanner(test.source, reporter)
			tokens := s.ScanTokens()

			p := parser.NewParser(tokens, reporter)
			statements, err := p.Parse()

			if err != nil {
				t.Fatalf("Parser returned nil for source: %s", test.source)
			}

			var out bytes.Buffer
			i := evaluation.NewInterpreter()
			i.Stdout = &out

			err = i.Interpret(statements)

			if test.expectedError {
				if err == nil {
					t.Errorf("Expected  error for [%s], but got none", test.source)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for [%s]: %v", test.source, err)
				}
				if i.LastValue != test.expected {
					t.Errorf("For [%s]: expected %v, got %v", test.source, test.expected, i.LastValue)
				}
			}
		})
	}
}
