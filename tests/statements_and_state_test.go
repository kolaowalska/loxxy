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

func TestStatementAndState(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedOut   string
		expectedError bool
	}{
		{"Print math", "print 1 + 2 * 3;", "7\n", false},
		{"Variable declaration and print", "var a = 1; print a;", "1\n", false},
		{"Uninitialized variable", "var a; print a;", "nil\n", false},
		{"Variables in expression", "var a = 1; var b = 2; print a + b;", "3\n", false},
		{"Variable reassignment", "var a = 1; a = 2; print a;", "2\n", false},
		{"Variable reassignment in print", "var a = 1; print a = 2;", "2\n", false},
		{"Right associativity of assignment", "var a; var b; a = b = 3; print a; print b;", "3\n3\n", false},
		{"Block scope variables", "var a = 1; { var b = 2; print a + b ; }", "3\n", false},
		{"Shadowing inner scope", "var a = 1; { var a = 2; print a; }", "2\n", false},
		{"Modifying outer scope variable in inner scope", "var a = 1; { a = 2; } print a;", "2\n", false},
		{"Deep nesting", "var a = \"global a\"; var b = \"global b\"; var c = \"global c\"; { var a = \"outer a\"; var b = \"outer b\"; { var a = \"inner a\"; print a; print b; print c; } print a; print b; print c; } print a; print b; print c;", "inner a\nouter b\nglobal c\nouter a\nouter b\nglobal c\nglobal a\nglobal b\nglobal c\n", false},
		{name: "Access undefined variable", source: "print undefinedVar;", expectedError: true},
		{name: "Invalid assignment target", source: "1 + 2 = 3;", expectedError: true},
		{name: "Missing semicolon after var", source: "var a = 1 print a;", expectedError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reporter := &testutils.TestReporter{}
			var out bytes.Buffer

			s := scanner.NewScanner(test.source, reporter)
			tokens := s.ScanTokens()
			if reporter.HadError && !test.expectedError {
				t.Fatalf("Scanner error on source: %s", test.source)
			}

			p := parser.NewParser(tokens, reporter)
			statements, err := p.Parse()

			if (reporter.HadError || err != nil) && !test.expectedError {
				t.Fatalf("Parser error on source: %s\n Error: %v", test.source, err)
			}
			if (reporter.HadError || err == nil) && test.expectedError {
				return
			}

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
