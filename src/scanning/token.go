package scanner

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
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
	return fmt.Sprintf("%v %s %v", t.TokenType, t.Lexeme, t.Literal)
}

func someTestFunction() {
	x := 42

}
