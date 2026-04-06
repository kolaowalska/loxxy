package parsing

import scanner "github.com/kolaowalska/loxxy/src/scanning"

type Parser struct {
	tokens   []scanner.Token
	current  int
	reporter scanner.ErrorReporter
}

func NewParser(tokens []scanner.Token, reporter scanner.ErrorReporter) *Parser {
	return &Parser{
		tokens:   tokens,
		reporter: reporter,
	}
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	return false
}

func (p *Parser) check(t scanner.TokenType) bool {
	return false
}

func (p *Parser) advance() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}
