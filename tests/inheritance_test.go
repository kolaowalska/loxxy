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

func TestInheritance(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      string
		expectedError bool
	}{
		{
			name: "INHERITANCE - Basic method inheritance",
			source: `
				class Doughnut {
					cook() { print "Fry until golden brown."; }
				}
				class BostonCream < Doughnut {}
				BostonCream().cook();
			`,
			expected:      "Fry until golden brown.\n",
			expectedError: false,
		},
		{
			name: "INHERITANCE - Calling super method",
			source: `
				class A {
					method() { print "A method"; }
				}
				class B < A {
					method() { print "B method"; super.method(); }
				}
				class C < B {}
				C().method();
			`,
			expected:      "B method\nA method\n",
			expectedError: false,
		},
		{
			name:          "INHERITANCE - Error on self-inheritance",
			source:        `class Oops < Oops {}`,
			expected:      "",
			expectedError: true, // Caught by Resolver
		},
		{
			name:          "INHERITANCE - Error inheriting from non-class",
			source:        `var NotAClass = "I am just a string"; class Subclass < NotAClass {}`,
			expected:      "",
			expectedError: true, // Caught by Interpreter at runtime
		},
		{
			name:          "INHERITANCE - Error using super outside class",
			source:        `print super.notEvenInAClass();`,
			expected:      "",
			expectedError: true, // Caught by Resolver
		},
		{
			name:          "INHERITANCE - Error using super in class with no parent",
			source:        `class Base { foo() { super.bar(); } }`,
			expected:      "",
			expectedError: true, // Caught by Resolver
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := &testutils.TestReporter{}

			// 1. scanning
			s := scanner.NewScanner(test.source, reporter)
			tokens := s.ScanTokens()

			// 2. parsing
			p := parser.NewParser(tokens, reporter)
			statements, err := p.Parse()

			if err != nil && !test.expectedError {
				t.Fatalf("Parser returned nil for source: %s\nError: %v", test.source, err)
			}

			// 3. resolving
			var out bytes.Buffer
			i := evaluation.NewInterpreter() // i := evaluation.NewInterpreter(&out)
			i.Stdout = &out

			resolver := resolving.NewResolver(i, reporter)
			_ = resolver.ResolveStatements(statements)

			if testutils.CheckError(t, test.expectedError, nil, reporter.HadError, "RESOLVING") {
				return
			}

			// 4. interpreting
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
