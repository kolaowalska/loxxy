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

//type MockReporter struct {
//	HadError bool
//}
//
//func (r *MockReporter) Error(line int, message string) {
//	r.HadError = true
//}
//
//func (r *MockReporter) TokenError(token scanner.Token, message string) {
//	r.HadError = true
//}

func TestStatementAndState(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectedOut string
		expectedErr bool
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
		{name: "Access undefined variable", source: "print undefinedVar;", expectedErr: true},
		{name: "Invalid assignment target", source: "1 + 2 = 3;", expectedErr: true},
		{name: "Missing semicolon after var", source: "var a = 1 print a;", expectedErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := &testutils.TestReporter{}
			var out bytes.Buffer

			s := scanner.NewScanner(tt.source, reporter)
			tokens := s.ScanTokens()
			if reporter.HadError && !tt.expectedErr {
				t.Fatalf("Scanner error on source: %s", tt.source)
			}

			p := parser.NewParser(tokens, reporter)
			statements, err := p.Parse()

			if (reporter.HadError || err != nil) && !tt.expectedErr {
				t.Fatalf("Parser error on source: %s\n Error: %v", tt.source, err)
			}
			if (reporter.HadError || err == nil) && tt.expectedErr {
				return
			}

			interpreter := evaluation.NewInterpreter()
			interpreter.Stdout = &out

			resolver := resolving.NewResolver(interpreter)
			err = resolver.ResolveStatements(statements)

			if err != nil {
				if tt.expectedErr {
					return
				}
				t.Fatalf("Resolver returned an error for source: %s\nError: %v", tt.source, err)
			}

			err = interpreter.Interpret(statements)

			if err != nil {
				if !tt.expectedErr {
					t.Fatalf("Unexpected error: %v", err)
				}
				return
			}

			if tt.expectedErr {
				t.Fatalf("Expected an error but got none for source: %s", tt.source)
			}

			if out.String() != tt.expectedOut {
				t.Errorf("\n Expected output: \n %q \n Actual output: \n %q", tt.expectedOut, out.String())
			}
		})
	}
}
