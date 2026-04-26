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

			resolver := resolving.NewResolver(i, reporter)
			_ = resolver.ResolveStatements(statements)

			if testutils.CheckError(t, test.expectedError, nil, reporter.HadError, "RESOLVING") {
				return
			}

			err = i.Interpret(statements)
			if testutils.CheckError(t, test.expectedError, err, reporter.HadError, "INTERPRETING") {
				return
			}
			if test.expectedError {
				t.Fatalf("expected an error for source: %s, but execution succeeded.", test.source)
			}
		})
	}
}
