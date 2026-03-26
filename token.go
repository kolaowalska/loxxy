package main

import "fmt"

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   any
	line      int
}

/*
	func newToken(_tokenType TokenType, _lexeme string, _literal any, _line int) *Token {
		return &Token{
			tokenType: _tokenType,
			lexeme:    _lexeme,
			literal:   _literal,
			line:      _line,
		}
	}
*/
func (t *Token) String() string {
	return fmt.Sprintf("%v %s %v", t.tokenType, t.lexeme, t.literal)
}
