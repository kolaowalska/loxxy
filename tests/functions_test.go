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
		{
			name:          "CALL - calling a non-function",
			source:        "var a = \"string\"; a();",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "CALL - arity mismatch (too many)",
			source:        "fun foo(a) {} foo(1, 2);",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "CALL - arity mismatch (too few)",
			source:        "fun foo(a, b) {} foo(1);",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "BUILDER - simple function with no return",
			source:        "fun sayHi(first, last) { print \"Hi, \" + first + \" \" + last + \"!\"; } sayHi(\"Dear\", \"Reader\");",
			expected:      "Hi, Dear Reader!\n",
			expectedError: false,
		},
		{
			name:          "BUILDER - parameters are locally scoped",
			source:        "fun foo(a) { print a; } foo(1); print a;",
			expected:      "1\n",
			expectedError: true,
		},
		{
			name:          "RETURN - standard return value",
			source:        "fun add(a, b) { return a + b; } print add(10, 20);",
			expected:      "30\n",
			expectedError: false,
		},
		{
			name:          "RETURN - early return skips rest of function",
			source:        "fun early() { return \"done\"; print \"never happens\"; } print early();",
			expected:      "done\n",
			expectedError: false,
		},
		{
			name:          "RETURN - empty return evaluates to nil",
			source:        "fun empty() { return; } print empty();",
			expected:      "nil\n",
			expectedError: false,
		},
		{
			name:          "RETURN - nested inside control flow",
			source:        "fun isEven(n) { if (n == 2) return true; return false; } print isEven(2); print isEven(3);",
			expected:      "true\nfalse\n",
			expectedError: false,
		},
		{
			name:          "INTEGRATION - recursive fibonacci",
			source:        "fun fib(n) { if (n <= 1) return n; return fib(n - 2) + fib(n - 1); } print fib(7);",
			expected:      "13\n",
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

			resolver := evaluation.NewResolver(i)
			err = resolver.ResolveStatements(statements)

			if err != nil {
				if test.expectedError {
					return
				}
				t.Fatalf("Resolver returned an error for source: %s\nError: %v", test.source, err)
			}

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
