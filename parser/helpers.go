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

func (p *parser) protoMethods() ([]ast.ProtoMethod, error) {
	methods := []ast.ProtoMethod{}
	method, err := p.protoMethod()

	if err != nil {
		return nil, err
	}

	methods = append(methods, method)

	for p.isThenEat(tokentype.COMMA) {
		method, err = p.protoMethod()

		if err != nil {
			return nil, err
		}

		methods = append(methods, method)
	}

	return methods, nil
}

func (p *parser) protoMethod() (ast.ProtoMethod, error) {
	k, err := p.fundamental()

	if err != nil {
		return ast.ProtoMethod{}, err
	}

	_, err = p.eat(tokentype.MINUS_GREAT)

	if err != nil {
		return ast.ProtoMethod{}, err
	}

	_, err = p.eat(tokentype.FN)

	if err != nil {
		return ast.ProtoMethod{}, err
	}

	m, err := p.function()

	if err != nil {
		return ast.ProtoMethod{}, err
	}

	return ast.ProtoMethod{K: k, M: m}, nil
}

type KeyVal = struct {
	Key ast.Expr
	Val ast.Expr
}

func (p *parser) recordKeyValue() (KeyVal, error) {
	k, err := p.fundamental()

	if err != nil {
		return KeyVal{}, err
	}

	_, err = p.eat(tokentype.MINUS_GREAT)

	if err != nil {
		return KeyVal{}, err
	}

	v, err := p.expression()

	if err != nil {
		return KeyVal{}, err
	}

	return KeyVal{Key: k, Val: v}, nil
}
