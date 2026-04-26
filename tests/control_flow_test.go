package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	"github.com/kolaowalska/loxxy/src/resolving"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func TestControlFlow(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
		{name: "IF - else test 1", source: "if (true) print \"good\"; else print \"bad\";", expected: "good\n", expectedError: false},
		{name: "IF - else test 2", source: "if (false) print \"bad\"; else print \"good\";", expected: "good\n", expectedError: false},
		{name: "IF - else block body", source: "if (false) nil; else { print \"block\"; }", expected: "block\n", expectedError: false},
		{name: "IF - dangling else 1", source: "if (true) if (false) print \"bad\"; else print \"good\";", expected: "good\n", expectedError: false},
		{name: "IF - dangling else 2", source: "if (false) if (true) print \"bad\"; else print \"good\";", expected: "", expectedError: false},
		{name: "IF - var in then", source: "if (true) var foo;", expected: "", expectedError: true},
		{name: "IF - var in else", source: "if (true) \"ok\"; else var foo;", expected: "", expectedError: true},
		{name: "IF - block body", source: "if (true) { print \"block\"; }", expected: "block\n", expectedError: false},
		{name: "IF - assignment in if", source: "var a = false; if (a = true) print a;", expected: "true\n", expectedError: false},
		{name: "IF - everything is true 1", source: "if (true) print true;", expected: "true\n", expectedError: false},
		{name: "IF - everything is true 2", source: "if (0) print 0;", expected: "0\n", expectedError: false},
		{name: "IF - everything is true 3", source: "if (\"\") print 0;", expected: "0\n", expectedError: false},
		{name: "IF - false is false", source: "if (false) print \"bad\"; else print \"false\";", expected: "false\n", expectedError: false},
		{name: "IF - nil is false", source: "if (nil) print \"bad\"; else print \"nil\";", expected: "nil\n", expectedError: false},
		{name: "AND - return first non-true argument 1", source: "print false and 1;", expected: "false\n", expectedError: false},
		{name: "AND - return first non-true argument 2", source: "print true and 1;", expected: "1\n", expectedError: false},
		{name: "AND - return first non-true argument 3", source: "print 1 and 2 and false;", expected: "false\n", expectedError: false},
		{name: "AND - return the last argument if all are true 1", source: "print 1 and true;", expected: "true\n", expectedError: false},
		{name: "AND - return the last argument if all are true 2", source: "print 1 and 2 and 3;", expected: "3\n", expectedError: false},
		{name: "AND - False is false", source: "print false and \"bad\";", expected: "false\n", expectedError: false},
		{name: "AND - Nil is false", source: "print nil and \"bad\";", expected: "nil\n", expectedError: false},
		{name: "AND - first true, then shouldn't matter 1", source: "print true and \"ok\";", expected: "ok\n", expectedError: false},
		{name: "AND - first true, then shouldn't matter 2", source: "print 0 and \"ok\";", expected: "ok\n", expectedError: false},
		{name: "AND - first true, then shouldn't matter 3", source: "print \"\" and \"ok\";", expected: "ok\n", expectedError: false},
		{name: "OR - return first true argument 1", source: "print 1 or true;", expected: "1\n", expectedError: false},
		{name: "OR - return first true argument 2", source: "print false or 1;", expected: "1\n", expectedError: false},
		{name: "OR - return first true argument 3", source: "print false or false or true;", expected: "true\n", expectedError: false},
		{name: "OR - return the last argument if all are false 1", source: "print false or false;", expected: "false\n", expectedError: false},
		{name: "OR - return the last argument if all are false 2", source: "print false or false or false;", expected: "false\n", expectedError: false},
		{name: "OR - False is false", source: "print false or \"ok\";", expected: "ok\n", expectedError: false},
		{name: "OR - Nil is false", source: "print nil or \"ok\";", expected: "ok\n", expectedError: false},
		{name: "OR - first true, then shouldn't matter 1", source: "print true or \"ok\";", expected: "true\n", expectedError: false},
		{name: "OR - first true, then shouldn't matter 2", source: "print 0 or \"ok\";", expected: "0\n", expectedError: false},
		{name: "OR - first true, then shouldn't matter 3", source: "print \"s\" or \"ok\";", expected: "s\n", expectedError: false},
		{
			name:          "Short circuit at the first true argument",
			source:        "var a = \"before\"; var b = \"before\"; (a = false) or (b = true) or (a = \"bad\"); print a; print b;",
			expected:      "false\ntrue\n",
			expectedError: false,
		},
		{
			name:          "Short circuit at the first false argument",
			source:        "var a = \"before\"; var b = \"before\"; (a = true) and (b = false) and (a = \"bad\"); print a; print b;",
			expected:      "true\nfalse\n",
			expectedError: false,
		},
		{
			name:          "if true",
			source:        "if (true) print \"ok\";",
			expected:      "ok\n",
			expectedError: false,
		},
		{
			name:          "if false",
			source:        "if (false) print \"ok\";",
			expected:      "",
			expectedError: false,
		},
		{
			name:          "while i",
			source:        "var i = 0; while (i < 4){print i; i = i + 1;}",
			expected:      "0\n1\n2\n3\n",
			expectedError: false,
		},
		{
			name:          "while no closed bracket",
			source:        "var i = 0; while (i < 4{print i; i = i + 1;}",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "WHILE - var in body",
			source:        "while (true) var foo;",
			expected:      "",
			expectedError: true,
		},
		{
			name:          "WHILE - without block body",
			source:        "var i = 0; while (i < 3) print i = i + 1;",
			expected:      "1\n2\n3\n",
			expectedError: false,
		},
		{
			name:          "for - print",
			source:        "for (var a = 1; a < 4; a = a + 1) {print a;}",
			expected:      "1\n2\n3\n",
			expectedError: false,
		},
		{
			name:          "for - no var",
			source:        "for (a = 1; a < 4; a = a + 1) {print a;}",
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "for - single statement body",
			source:        "for (var i = 0; i < 3; i = i + 1) print i;",
			expected:      "0\n1\n2\n",
			expectedError: false,
		},
		{
			name:          "for - missing initializer",
			source:        "var i = 0; for (; i < 3; i = i + 1) print i;",
			expected:      "0\n1\n2\n",
			expectedError: false,
		},
		{
			name:          "for - missing increment",
			source:        "for (var i = 0; i < 3;) { print i; i = i + 1; }",
			expected:      "0\n1\n2\n",
			expectedError: false,
		},
		{
			name:          "for - missing initializer and increment",
			source:        "var i = 0; for (; i < 3;) { print i; i = i + 1; }",
			expected:      "0\n1\n2\n",
			expectedError: false,
		},
		{
			name:          "for - expression statement initializer",
			source:        "var a; for (a = 0; a < 3; a = a + 1) print a;",
			expected:      "0\n1\n2\n",
			expectedError: false,
		},
		{
			name:          "for - variable scoping",
			source:        "for (var i = 0; i < 2; i = i + 1) { print i; } print i;",
			expected:      "",
			expectedError: true, // should trigger a runtime error
		},
		{
			name:          "for - variable shadowing",
			source:        "var i = 5; for (var i = 0; i < 2; i = i + 1) { print i; } print i;",
			expected:      "0\n1\n5\n",
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
					t.Errorf("Expected error for [%s], but got none", test.source)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for [%s]: %v", test.source, err)
				}
				if out.String() != test.expected {
					t.Errorf("For source:\n%s\n\nExpected:\n%v\n\nGot:\n%v", test.source, test.expected, out.String())
				}
			}
		})
	}
}
