package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func TestFunctions(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{
			name:          "CALL - native clock function",
			source:        "var t = clock(); print t > 0;",
			expected:      "true\n",
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := TestReporter{}
			s := scanner.NewScanner(test.source, reporter)
			tokens := s.ScanTokens()

			p := parser.NewParser(tokens, reporter)
			statements, err := p.Parse()

			if err != nil {
				if test.expectedError {
					return
				}
				t.Fatalf("Parser returned nil for source: %s\nError: %v", test.source, err)
			}

			var out bytes.Buffer
			i := evaluation.NewInterpreter() // i := evaluation.NewInterpreter(&out)
			i.Stdout = &out
			err = i.Interpret(statements)

			if err != nil {
				if test.expectedError {
					return
				}
				t.Fatalf("Interpreter returned an error for source: %s\nError: %v", test.source, err)
			}
			if test.expectedError {
				t.Fatalf("Expected an error for source: %s, but execution succeeded.", test.source)
			}
			if out.String() != test.expected {
				t.Errorf("For source:\n%s\n\nExpected:\n%v\n\nGot:\n%v", test.source, test.expected, out.String())
			}
		})
	}
}
