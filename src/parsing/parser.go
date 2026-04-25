package parsing

import (
	"fmt"

	"github.com/kolaowalska/loxxy/src/representation"
	scanner "github.com/kolaowalska/loxxy/src/scanning"
)

const (
	msgSemicolon         = "expect ';' after variable declaration"
	msgExpression        = "expect expression"
	msgVariableName      = "expect variable name"
	msgRightParen        = "expect ')' after expression"
	msgInvalidAssignment = "invalid assignment target"
	msgRightCurlyParen   = "expect '}' after block."
)

type ErrorReporter interface {
	Error(line int, message string)
	TokenError(t scanner.Token, message string)
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

func (p *Parser) Parse() ([]representation.Stmt, error) {
	var statements []representation.Stmt

	for !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, dec)
	}
	return statements, nil
}

func (p *Parser) expression() (representation.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (representation.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*representation.Variable); ok {
			name := v.Name
			return &representation.Assign{Name: name, Value: value}, nil
		}

		_ = p.error(equals, msgInvalidAssignment)
	}
	return expr, nil
}

func (p *Parser) or() (representation.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &representation.Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) and() (representation.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &representation.Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) declaration() (representation.Stmt, error) {
	if p.match(scanner.VAR) {
		return p.varDeclaration()
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		// NOTE: we do not return err (like in a book). change it if necessary
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) varDeclaration() (representation.Stmt, error) {
	name, err := p.consume(scanner.IDENTIFIER, msgVariableName)
	if err != nil {
		return nil, err
	}

	var initializer representation.Expr
	if p.match(scanner.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(scanner.SEMICOLON, msgSemicolon)
	if err != nil {
		return nil, err
	}

	return &representation.Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (representation.Stmt, error) {
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.FOR) {
		return p.forStatement()
	}
	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return &representation.Block{Statements: block}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() (representation.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "expect '(' after 'if'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after if condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch representation.Stmt
	if p.match(scanner.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &representation.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) printStatement() (representation.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMICOLON, msgSemicolon)
	if err != nil {
		return nil, err
	}
	return &representation.Print{Expression: value}, nil
}

func (p *Parser) expressionStatement() (representation.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.SEMICOLON, msgSemicolon)
	if err != nil {
		return nil, err
	}
	return &representation.Expression{Expression: expr}, nil
}

func (p *Parser) block() ([]representation.Stmt, error) {
	var statements []representation.Stmt

	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		dec, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, dec)
	}
	_, err := p.consume(scanner.RIGHT_BRACE, msgRightCurlyParen)
	if err != nil {
		return nil, err
	}
	return statements, nil
}

// control flow ----------------------------------------------

func (p *Parser) forStatement() (representation.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "expect '(' after 'for'")
	if err != nil {
		return nil, err
	}

	var initializer representation.Stmt
	if p.match(scanner.SEMICOLON) {
		initializer = nil
	} else if p.match(scanner.VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var condition representation.Expr = nil
	if !p.check(scanner.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.SEMICOLON, "expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var increment representation.Expr = nil
	if !p.check(scanner.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after for clauses")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	// desugaring
	if increment != nil {
		body = &representation.Block{
			Statements: []representation.Stmt{body, &representation.Expression{Expression: increment}},
		}
	}

	if condition == nil {
		condition = &representation.Literal{Value: true}
	}

	body = &representation.While{Condition: condition, Body: body}

	if initializer != nil {
		body = &representation.Block{
			Statements: []representation.Stmt{initializer, body},
		}
	}

	return body, nil
}

// -----------------------------------------------------------

func (p *Parser) whileStatement() (representation.Stmt, error) {
	_, err := p.consume(scanner.LEFT_PAREN, "expect '(' after 'while'")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(scanner.RIGHT_PAREN, "expect ')' after condition")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &representation.While{Condition: condition, Body: body}, nil

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

func (p *Parser) consume(t scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return p.peek(), p.error(p.peek(), message)
}

func (p *Parser) error(t scanner.Token, message string) error {
	p.reporter.TokenError(t, message)
	return fmt.Errorf("%s", message)
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

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
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
	return p.call()
}

func (p *Parser) call() (representation.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(scanner.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee representation.Expr) (representation.Expr, error) {
	var arguments []representation.Expr

	if !p.check(scanner.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				_ = p.error(p.peek(), "can't have more than 255 arguments.")
			}
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if !p.match(scanner.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(scanner.RIGHT_PAREN, "expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &representation.Call{Callee: callee, Paren: paren, Args: arguments}, nil
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
	if p.match(scanner.IDENTIFIER) {
		return &representation.Variable{Name: p.previous()}, nil
	}
	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(scanner.RIGHT_PAREN, msgRightParen)
		if err != nil {
			return nil, err
		}

		return &representation.Grouping{Expression: expr}, nil
	}

	return nil, p.error(p.peek(), msgExpression)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == scanner.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case scanner.CLASS, scanner.FUN, scanner.VAR, scanner.FOR,
			scanner.IF, scanner.WHILE, scanner.PRINT, scanner.RETURN:
			return
		default:
			_ = fmt.Errorf("it's not supposed to go there, error in func synchronize")
		}
		p.advance()
	}
}
