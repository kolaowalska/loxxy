package scanner

import (
	"github.com/kolaowalska/loxxy/src/reports"
	"testing"
)

func TestScanner_ValidTokens(t *testing.T) {
	defer reporter.Clear()

	tests := []struct {
		name     string
		source   string
		expected []TokenType
	}{
		{"single characters", "(){},.-+;*/", []TokenType{
			LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE,
			COMMA, DOT, MINUS, PLUS, SEMICOLON, STAR, SLASH, EOF,
		}},
		{"operators", "!= == <= >= = < >", []TokenType{
			BANG_EQUAL, EQUAL_EQUAL, LESS_EQUAL, GREATER_EQUAL,
			EQUAL, LESS, GREATER, EOF,
		}},
		{"strings", "\"hello\" \"ala has a cat\"", []TokenType{
			STRING, STRING, EOF,
		}},
		{"numbers", "123.45 123 .456 123.", []TokenType{
			NUMBER, NUMBER, DOT, NUMBER, NUMBER, DOT, EOF,
		}},
		{"keywords and identifiers", "var orchid = 10;", []TokenType{
			VAR, IDENTIFIER, EQUAL, NUMBER, SEMICOLON, EOF,
		}},
		{"keywords", "and class else false for fun if nil or print return super this true var while", []TokenType{
			AND, CLASS, ELSE, FALSE, FOR, FUN, IF, NIL, OR, PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE, EOF,
		}},
		{"comments", "// comments about code", []TokenType{
			EOF,
		}},
		{"whitespaces", "a b\rc\td\ne", []TokenType{
			IDENTIFIER, IDENTIFIER, IDENTIFIER, IDENTIFIER, IDENTIFIER, EOF,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter.Clear()
			scanner := NewScanner(tt.source)
			tokens := scanner.ScanTokens()

			if reporter.HadError {
				t.Fatalf("reporter wywalil error a nie powinien")
			}

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, expectedType := range tt.expected {
				if tokens[i].TokenType != expectedType {
					t.Errorf("token %d: expected %v, got %v instead", i, expectedType, tokens[i].TokenType)
				}
			}
		})
	}
}

func TestScanner_LexicalErrors(t *testing.T) {
	defer reporter.Clear()

	tests := []struct {
		name   string
		source string
	}{
		{"unexpected character.", "@"},
		{"unexpected character.", "#"},
		{"unexpected character.", "^"},
		{"unterminated string.", "\"this string has no end"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter.Clear() // Reset before scanning
			scanner := NewScanner(tt.source)
			scanner.ScanTokens()

			// ERROR IS EXPECTED
			if !reporter.HadError {
				t.Errorf("expected scanner to report an error for %q, but it didn't", tt.source)
			}
		})
	}
}
