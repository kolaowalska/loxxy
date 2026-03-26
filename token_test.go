package main

import "testing"

func TestTokenString(t *testing.T) {
	tests := []struct {
		name     string
		token    *Token
		expected string
	}{
		{
			name:     "Simple symbol",
			token:    &Token{tokenType: LEFT_PAREN, lexeme: "(", literal: nil, line: 1},
			expected: "LEFT_PAREN ( <nil>",
		},
		{
			name:     "Number Literal",
			token:    &Token{NUMBER, "123.45", 123.45, 1},
			expected: "NUMBER 123.45 123.45",
		},
		{
			name:     "String Literal",
			token:    &Token{STRING, "\"hello\"", "hello", 2},
			expected: "STRING \"hello\" hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.token.String()
			if actual != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
