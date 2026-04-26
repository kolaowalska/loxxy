package testutils

import (
	"fmt"

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
