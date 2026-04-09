package parsing

import (
	"fmt"

	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

type ErrorReporter interface {
	Error(line int, message string)
}

type Parser struct {
	tokens   []scanner.Token
	current  int
	reporter ErrorReporter
}

func NewParser(tokens []scanner.Token, reporter ErrorReporter) *Parser {
	return &Parser{
		tokens:   tokens,
		current:  0,
		reporter: reporter,
	}
}

func (p *Parser) Parse() representation.Expr {
	expr, err := p.expression()
	if err != nil {
		return nil
	}
	return expr
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

// TODO: empty body,
func (p *Parser) consume(t scanner.TokenType, message string) {
	return
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

func (p *Parser) equality() (representation.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		operator := p.previous()

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &representation.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

// comparison: term (( ">" | ">=" | "<" | "<=" ) term )*
func (p *Parser) comparison() (representation.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &representation.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// term: factor (( "-" | "+" ) factor )*
func (p *Parser) term() (representation.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.MINUS, scanner.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &representation.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// factor: unary (( "/" | "*" ) unary )*
func (p *Parser) factor() (representation.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.SLASH, scanner.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &representation.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

// unary -> ( "!" | "-" ) unary | primary
func (p *Parser) unary() (representation.Expr, error) {
	if p.match(scanner.BANG, scanner.MINUS) {
		operator := p.previous()
		right, err := p.unary() // JESLI COS SIE PETLI TO PEWNIE DLATEGO XD
		if err != nil {
			return nil, err
		}

		return &representation.Unary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

// primary: NUMBER | STRING | "true" | "false" | "nil" | "(" expr ")" | IDENTIFIER
func (p *Parser) primary() (representation.Expr, error) {
	if p.match(scanner.TRUE) {
		return &representation.Literal{Value: true}, nil
	}
	if p.match(scanner.FALSE) {
		return &representation.Literal{Value: false}, nil
	}
	if p.match(scanner.NIL) {
		return &representation.Literal{Value: nil}, nil
	}
	if p.match(scanner.NUMBER, scanner.STRING) {
		return &representation.Literal{Value: p.previous().Literal}, nil
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		p.consume(scanner.RIGHT_PAREN, "Expect '?' after expression.")
		return &representation.Grouping{Expression: expr}, nil
	}

	p.reporter.Error(p.peek().Line, "Expect expression.")

	return &representation.Literal{Value: p.previous().Literal}, fmt.Errorf("Expect expression.", p.peek().Line)
}
