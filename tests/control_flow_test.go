package tests

import (
	"bytes"
	"testing"

	"github.com/kolaowalska/loxxy/src/evaluation"
	parser "github.com/kolaowalska/loxxy/src/parsing"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

func TestControlFlow(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expected      any
		expectedError bool
	}{
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
		// {
		// 	name:          "while no closed bracket",
		// 	source:        "var i = 0; while (i < 4{print i; i = i + 1;}",
		// 	expected:      nil,
		// 	expectedError: true,
		// },
		{
			name:          "for print",
			source:        "for (var a = 1; a < 4; a = a + 1) {print a;}",
			expected:      "1\n2\n3\n",
			expectedError: false,
		},
		{
			name:          "for no var",
			source:        "for (a = 1; a < 4; a = a + 1) {print a;}",
			expected:      nil,
			expectedError: true,
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
				t.Fatalf("Parser returned nil for source: %s\nError: %v", test.source, err)
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
				if out.String() != test.expected {
					t.Errorf("For source:\n%s\n\nExpected:\n%v\n\nGot:\n%v", test.source, test.expected, out.String())
				}
			}
		})
	}
}
