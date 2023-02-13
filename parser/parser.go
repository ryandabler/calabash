package parser

import (
	"calabash/ast"
	"calabash/errors"
	"calabash/internal/tokentype"
	"calabash/lexer/tokens"
	"fmt"
)

type parser struct {
	tokens []tokens.Token
	i      int
}

func (p *parser) Parse() ([]ast.Node, error) {
	return p.program()
}

func (p *parser) atEnd() bool {
	return p.tokens[p.i].Type == tokentype.EOF
}

func (p *parser) eat(ts ...tokentype.Tokentype) (tokens.Token, error) {
	if p.atEnd() {
		return tokens.Token{}, errors.ParseError{Msg: "Unexpected end of input"}
	}

	for _, v := range ts {
		if v == p.tokens[p.i].Type {
			t := p.tokens[p.i]
			p.next()

			return t, nil
		}
	}

	e := errors.ParseError{Msg: fmt.Sprintf("Token of type %q did match types %#v", p.tokens[p.i].Type, ts)}
	return tokens.Token{}, e
}

func (p *parser) is(ts ...tokentype.Tokentype) bool {
	if p.atEnd() {
		return false
	}

	for _, v := range ts {
		if v == p.tokens[p.i].Type {
			return true
		}
	}

	return false
}

func (p *parser) isThenEat(ts ...tokentype.Tokentype) bool {
	is := p.is(ts...)

	if is {
		p.eat(ts...)
	}

	return is
}

func (p *parser) next() {
	p.i++
}

func (p *parser) current() tokens.Token {
	return p.tokens[p.i]
}

func (p *parser) program() ([]ast.Node, error) {
	ts := []ast.Node{}

	for !p.atEnd() {
		if p.isThenEat(tokentype.LET) {
			n, err := p.variableDecl()

			if err != nil {
				return []ast.Node{}, err
			}

			ts = append(ts, n)
			continue
		}

		expr, err := p.expression()

		if err != nil {
			return []ast.Node{}, err
		}

		// If the next token is a comma or equals sign, we are processing an
		// assignment statement and not an expression.
		if p.is(tokentype.COMMA, tokentype.EQUAL) {
			n, err := p.assignment(expr)

			if err != nil {
				return nil, err
			}

			ts = append(ts, n)
			continue
		}

		ts = append(ts, expr)
	}

	return ts, nil
}

func (p *parser) variableDecl() (ast.Node, error) {
	names, err := p.varDeclarationNames()

	if err != nil {
		return nil, err
	}

	// No initializers are specified for this assignment
	if p.isThenEat(tokentype.SEMICOLON) {
		return ast.VarDeclStmt{Names: names, Values: []ast.Expr{}}, nil
	}

	// Gather initializing values
	_, err = p.eat(tokentype.EQUAL)

	if err != nil {
		return nil, err
	}

	inits, err := p.commaExpressions()

	if err != nil {
		return nil, err
	}

	_, err = p.eat(tokentype.SEMICOLON)

	if err != nil {
		n := p.current()
		return nil, errors.ParseError{Msg: fmt.Sprintf("Missing semicolon at %d:%d", n.Position.Row, n.Position.Col)}
	}

	return ast.VarDeclStmt{Names: names, Values: inits}, nil
}

func (p *parser) assignment(fst ast.Expr) (ast.Node, error) {
	ident, ok := fst.(ast.IdentifierExpr)

	if !ok {
		return nil, errors.ParseError{Msg: "Expected identifier for first element of assignment statement"}
	}

	ns := []tokens.Token{ident.Name}

	// Obtain the list of identifiers being reassigned
	for p.isThenEat(tokentype.COMMA) {
		n, err := p.eat(tokentype.IDENTIFIER)

		if err != nil {
			return nil, err
		}

		ns = append(ns, n)
	}

	// Skip past the equals sign
	_, err := p.eat(tokentype.EQUAL)

	if err != nil {
		return nil, err
	}

	// Get the list of expressions being bound to the above identifiers
	exprs, err := p.commaExpressions()

	if err != nil {
		return nil, err
	}

	// Skip semicolon
	_, err = p.eat(tokentype.SEMICOLON)

	if err != nil {
		return nil, err
	}

	return ast.AssignmentStmt{Names: ns, Values: exprs}, nil
}

func (p *parser) expression() (ast.Expr, error) {
	return p.booleanOr()
}

func (p *parser) booleanOr() (ast.Expr, error) {
	left, err := p.booleanAnd()

	if err != nil {
		return nil, err
	}

	for p.is(tokentype.STROKE_STROKE) {
		op, _ := p.eat(tokentype.STROKE_STROKE)
		right, err := p.booleanAnd()

		if err != nil {
			return nil, err
		}

		left = ast.BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	return left, nil
}

func (p *parser) booleanAnd() (ast.Expr, error) {
	left, err := p.equality()

	if err != nil {
		return nil, err
	}

	for p.is(tokentype.AMPERSAND_AMPERSAND) {
		op, _ := p.eat(tokentype.AMPERSAND_AMPERSAND)
		right, err := p.equality()

		if err != nil {
			return nil, err
		}

		left = ast.BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: op,
		}
	}

	return left, nil
}

func (p *parser) equality() (ast.Expr, error) {
	left, err := p.comparison()

	if err != nil {
		return nil, err
	}

	if p.is(equalityTokens...) {
		op, _ := p.eat(equalityTokens...)
		right, err := p.comparison()

		if err != nil {
			return nil, err
		}

		return ast.BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: op,
		}, nil
	}

	return left, nil
}

func (p *parser) comparison() (ast.Expr, error) {
	left, err := p.addition()

	if err != nil {
		return nil, err
	}

	if p.is(comparisonTokens...) {
		op, _ := p.eat(comparisonTokens...)
		right, err := p.addition()

		if err != nil {
			return nil, err
		}

		return ast.BinaryExpr{
			Left:     left,
			Right:    right,
			Operator: op,
		}, nil
	}

	return left, nil
}

func (p *parser) addition() (ast.Expr, error) {
	l, err := p.multiplication()

	if err != nil {
		return nil, err
	}

	for p.is(tokentype.PLUS, tokentype.MINUS) {
		op, _ := p.eat(tokentype.PLUS, tokentype.MINUS)
		r, err := p.multiplication()

		if err != nil {
			return nil, err
		}

		l = ast.BinaryExpr{
			Left:     l,
			Right:    r,
			Operator: op,
		}
	}

	return l, nil
}

func (p *parser) multiplication() (ast.Expr, error) {
	l, err := p.exponentiation()

	if err != nil {
		return nil, err
	}

	for p.is(tokentype.ASTERISK, tokentype.SLASH) {
		op, _ := p.eat(tokentype.ASTERISK, tokentype.SLASH)
		r, err := p.exponentiation()

		if err != nil {
			return nil, err
		}

		l = ast.BinaryExpr{
			Left:     l,
			Right:    r,
			Operator: op,
		}
	}

	return l, nil
}

func (p *parser) exponentiation() (ast.Expr, error) {
	l, err := p.unary()

	if err != nil {
		return nil, err
	}

	for p.is(tokentype.ASTERISK_ASTERISK) {
		op, _ := p.eat(tokentype.ASTERISK_ASTERISK)
		r, err := p.unary()

		if err != nil {
			return nil, err
		}

		l = ast.BinaryExpr{
			Left:     l,
			Right:    r,
			Operator: op,
		}
	}

	return l, nil
}

func (p *parser) unary() (ast.Expr, error) {
	if p.is(tokentype.MINUS) {
		op, _ := p.eat(tokentype.MINUS)
		expr, err := p.unary()

		if err != nil {
			return nil, err
		}

		return ast.UnaryExpr{Operator: op, Expr: expr}, nil
	}

	return p.fundamental()
}

func (p *parser) fundamental() (ast.Expr, error) {
	if p.atEnd() {
		return nil, errors.ParseError{Msg: "Unexpected end of input"}
	}

	if p.isThenEat(tokentype.LEFT_PAREN) {
		e, err := p.expression()

		if err != nil {
			return nil, err
		}

		_, err = p.eat(tokentype.RIGHT_PAREN)

		if err != nil {
			return nil, err
		}

		return ast.GroupingExpr{Expr: e}, nil
	}

	if p.is(tokentype.NUMBER) {
		n, _ := p.eat(tokentype.NUMBER)
		return ast.NumericLiteralExpr{Value: n}, nil
	}

	if p.is(tokentype.STRING) {
		s, _ := p.eat(tokentype.STRING)
		return ast.StringLiteralExpr{Value: s}, nil
	}

	if p.is(tokentype.BOTTOM) {
		s, _ := p.eat(tokentype.BOTTOM)
		return ast.BottomLiteralExpr{Token: s}, nil
	}

	if p.is(tokentype.IDENTIFIER) {
		s, _ := p.eat(tokentype.IDENTIFIER)
		return ast.IdentifierExpr{Name: s}, nil
	}

	if p.is(tokentype.TRUE, tokentype.FALSE) {
		b, _ := p.eat(tokentype.TRUE, tokentype.FALSE)
		return ast.BooleanLiteralExpr{Value: b}, nil
	}

	t := p.tokens[p.i]
	return nil, errors.ParseError{Msg: fmt.Sprintf("Malformed expression at %d: %d", t.Position.Row, t.Position.Col)}
}

func New(ts []tokens.Token) *parser {
	return &parser{tokens: ts}
}
