package parser

import (
	"calabash/ast"
	"calabash/internal/tokentype"
)

func (p *parser) varDeclarationNames() ([]ast.Identifier, error) {
	ns := []ast.Identifier{}
	n, err := p.varName()

	if err != nil {
		return nil, err
	}

	ns = append(ns, n)

	for p.isThenEat(tokentype.COMMA) {
		n, err = p.varName()

		if err != nil {
			return nil, err
		}

		ns = append(ns, n)
	}

	return ns, nil
}

func (p *parser) varName() (ast.Identifier, error) {
	i := ast.Identifier{}

	if p.isThenEat(tokentype.MUT) {
		i.Mut = true
	}

	ident, err := p.eat(tokentype.IDENTIFIER)

	if err != nil {
		return i, err
	}

	i.Name = ident

	return i, nil
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
