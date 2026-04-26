package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	"github.com/kolaowalska/loxxy/src/resolving"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
	"github.com/kolaowalska/loxxy/src/testutils"
)

//type TestReporter struct{}
//
//func (r TestReporter) Error(line int, message string)                 {}
//func (r TestReporter) TokenError(token scanner.Token, message string) {}

func TestPipeline(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{"Simple math", "print 1 + 2 * 3;", "7\n", false},
		{"Grouping", "print (1 + 2) * 3;", "9\n", false},
		{"String concat", "print \"ala\" + \"kot\";", "alakot\n", false},
		{"Comparison and booleans", "print 10 >= 5 == true;", "true\n", false},
		{"Equality mixed", "print \"ala\" == 123;", "false\n", false},
		{"Complex test 1", "print (-1 + 5) * 2 / 4;", "2\n", false},
		{"Complex test 2", "print (5 - (3 - 1)) + -1;", "2\n", false},

		{"Subtract string from number", "\"ala\" - 123;", nil, true},
		{"Unary minus on string", "-\"ala\";", nil, true},
		{"Add booleans", "true + false;", nil, true},
		{"Divide by zero", "10 / 0;", nil, true},
		{"Greater comparison on strings", "\"a\" > \"b\";", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := &testutils.TestReporter{}

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

			resolver := resolving.NewResolver(i)
			err = resolver.ResolveStatements(statements)

			if err != nil {
				if test.expectedError {
					return
				}
				t.Fatalf("Resolver returned an error for source: %s\nError: %v", test.source, err)
			}

			err = i.Interpret(statements)

			if test.expectedError {
				if err == nil {
					t.Errorf("Expected  error for [%s], but got none", test.source)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for [%s]: %v", test.source, err)
				}
				if out.String() != test.expected {
					t.Errorf("For [%s]: expected %v, got %v", test.source, test.expected, out.String())
				}
			}
		})
	}
}
