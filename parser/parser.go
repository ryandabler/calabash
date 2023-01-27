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

func (p *parser) next() {
	p.i++
}

func (p *parser) program() ([]ast.Node, error) {
	ts := []ast.Node{}

	for !p.atEnd() {
		expr, err := p.expression()

		if err != nil {
			return []ast.Node{}, err
		}

		ts = append(ts, expr)
	}

	return ts, nil
}

func (p *parser) expression() (ast.Expr, error) {
	return p.addition()
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

	if p.is(tokentype.LEFT_PAREN) {
		p.eat(tokentype.LEFT_PAREN)
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

	t := p.tokens[p.i]
	return nil, errors.ParseError{Msg: fmt.Sprintf("Malformed expression at %d: %d", t.Position.Row, t.Position.Col)}
}

func New(ts []tokens.Token) *parser {
	return &parser{tokens: ts}
}
