package testutils

import (
	"fmt"
	"testing"

	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type TestReporter struct {
	HadError    bool
	LastMessage string
}

func (r *TestReporter) Error(line int, message string) {
	r.HadError = true
	r.LastMessage = message
	fmt.Printf("[line %d] error: %s\n", line, message)
}

func (r *TestReporter) TokenError(t scanner.Token, message string) {
	r.HadError = true
	r.LastMessage = message
	fmt.Printf("[line %d] error at '%s': %s\n", t.Line, t.Lexeme, message)
}

func CheckError(t *testing.T, expectedError bool, err error, hadError bool, phase string) bool {
	t.Helper() // marking this as a helper so error line numbers point to the test file

	if err != nil || hadError {
		if expectedError {
			return true
		}
		t.Fatalf("unexpected error during %s phase. error: %v", phase, err)
	}

	return false
}
