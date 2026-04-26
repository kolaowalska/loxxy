package resolving_test

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	"github.com/kolaowalska/loxxy/src/resolving"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
	"github.com/kolaowalska/loxxy/src/testutils"
)

type TestReporter struct{}

func (r TestReporter) Error(line int, message string)             {}
func (r TestReporter) TokenError(t scanner.Token, message string) {}

func TestResolvingAndBinding(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{
			name: "RESOLVER - closures capture the correct environment depth",
			source: `
				var a = "global";
				{
					fun showA() {
						print a;
					}

					showA();
					var a = "block";
					showA();
				}
			`,
			expected:      "global\nglobal\n",
			expectedError: false,
		},
		{
			name: "RESOLVER ERROR - cannot read local variable in its own initializer",
			source: `
				var a = "outer";
				{
					var a = a; 
				}
			`,
			expected:      "",
			expectedError: true, // should trigger: "can't read local variable in its own initializer"
		},
		{
			name: "RESOLVER ERROR - cannot declare two variables with the same name in the same local scope",
			source: `
				{
					var a = "first";
					var a = "second";
				}
			`,
			expected:      "",
			expectedError: true, // should trigger: "already a variable with this name in this scope"
		},
		{
			name: "RESOLVER - global variable redeclaration is allowed",
			source: `
				var a = "first";
				var a = "second";
				print a;
			`,
			expected:      "second\n",
			expectedError: false, // lox allows redefining globals, so the resolver should not error here!
		},
		{
			name: "RESOLVER - assignment resolves to the correct shadowed scope",
			source: `
				var a = "global";
				{
					var a = "local";
					a = "modified local";
					print a;
				}
				print a;
			`,
			expected:      "modified local\nglobal\n",
			expectedError: false,
		},
		{
			name: "RESOLVER - deeply nested variable resolution",
			source: `
				var a = "global";
				{
					var b = "outer";
					{
						var c = "inner";
						print a;
						print b;
						print c;
					}
				}
			`,
			expected:      "global\nouter\ninner\n",
			expectedError: false,
		},
		{
			name: "RESOLVER - deeply nested shadowing resolves correctly",
			source: `
				var a = "global";
				{
					var a = "outer";
					{
						var a = "inner";
						print a;
					}
					print a;
				}
				print a;
			`,
			expected:      "inner\nouter\nglobal\n",
			expectedError: false,
		},
		{
			name: "RESOLVER - assignment respects lexical scope",
			source: `
				var a = "global";
				{
					fun assign() {
						a = "assigned";
					}
					var a = "inner";
					assign();
					print a;
				}
				print a;
			`,
			expected:      "inner\nassigned\n",
			expectedError: false,
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
