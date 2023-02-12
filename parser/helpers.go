package parser

import (
	"calabash/ast"
	"calabash/internal/tokentype"
	"calabash/lexer/tokens"
)

func (p *parser) assignmentNames() ([]tokens.Token, error) {
	ident, err := p.eat(tokentype.IDENTIFIER)
	ns := []tokens.Token{}

	if err != nil {
		return ns, err
	}

	ns = append(ns, ident)

	for p.isThenEat(tokentype.COMMA) {
		ident, err := p.eat(tokentype.IDENTIFIER)

		if err != nil {
			return []tokens.Token{}, err
		}

		ns = append(ns, ident)
	}

	return ns, nil
}

func (p *parser) commaExpressions() ([]ast.Expr, error) {
	e, err := p.expression()
	es := []ast.Expr{}

	if err != nil {
		return es, err
	}

	es = append(es, e)

	for p.isThenEat(tokentype.COMMA) {
		e, err := p.expression()

		if err != nil {
			return []ast.Expr{}, err
		}

		es = append(es, e)
	}

	return es, nil
}
