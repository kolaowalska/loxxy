package parsing

import (
	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) expression() (representation.Expr, error) {
	return p.equality()
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == scanner.EOF
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

// TODO
func (p *Parser) equality() (representation.Expr, error) {
	return nil, nil
}

func (p *Parser) comparison() (representation.Expr, error) {
	return nil, nil
}

func (p *Parser) term() (representation.Expr, error) {
	return nil, nil
}

func (p *Parser) factor() (representation.Expr, error) {
	return nil, nil
}

func (p *Parser) unary() (representation.Expr, error) {
	return nil, nil
}

func (p *Parser) primary() (representation.Expr, error) {
	return nil, nil
}
