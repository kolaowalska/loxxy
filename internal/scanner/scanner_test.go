package scanner

import (
	"testing"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		name string 
		source string 
		expected []TokenType
	}{
		{"single characters", "(){}", []TokenType{LEFT_PARENTHESIS, RIGHT_PARENTHESIS, LEFT_BRACE, RIGHT_BRACE, EOF}},
		{"operators", "!= ==", []TokenType{NOT_EQUAL, EQUAL_EQUAL, EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := NewScanner(tt.source)
			tokens := scanner.ScanTokens()

			if len(tokens) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}
		})
	}
}
