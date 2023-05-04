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
		t, err := p.stmtOrExpr()

		if err != nil {
			return []ast.Node{}, err
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func (p *parser) stmtOrExpr() (ast.Node, error) {
	if p.isThenEat(tokentype.LET) {
		n, err := p.variableDecl()

		if err != nil {
			return nil, err
		}

		return n, nil
	}

	if p.isThenEat(tokentype.IF) {
		n, err := p.ifStmt()

		if err != nil {
			return nil, err
		}

		return n, nil
	}

	if p.isThenEat(tokentype.RETURN) {
		n, err := p.retStmt()

		if err != nil {
			return nil, err
		}

		return n, nil
	}

	expr, err := p.expression()

	if err != nil {
		return nil, err
	}

	// If the next token is a comma or equals sign, we are processing an
	// assignment statement and not an expression.
	if p.is(tokentype.COMMA, tokentype.EQUAL) {
		n, err := p.assignment(expr)

		if err != nil {
			return nil, err
		}

		return n, nil
	}

	return expr, nil
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

func (p *parser) ifStmt() (ast.Node, error) {
	var varDecl ast.Node
	var decls ast.VarDeclStmt
	var err error

	// Obtain option variable declarations
	if p.isThenEat(tokentype.LET) {
		varDecl, err = p.variableDecl()

		if err != nil {
			return nil, err
		}

		decls, _ = varDecl.(ast.VarDeclStmt)
	}

	// Obtain boolean condition
	cond, err := p.expression()

	if err != nil {
		return nil, err
	}

	// Obtain `then` portion of statement
	then, err := p.blockStmt()

	if err != nil {
		return nil, err
	}

	// Obtain optional `else` portion
	var elseBlk ast.Node

	if p.isThenEat(tokentype.ELSE) {
		if p.isThenEat(tokentype.IF) {
			elseBlk, err = p.ifStmt()

			if err != nil {
				return nil, err
			}
		} else {
			elseBlk, err = p.blockStmt()

			if err != nil {
				return nil, err
			}

		}
	}

	return ast.IfStmt{Decls: decls, Condition: cond, Then: then, Else: elseBlk}, nil
}

func (p *parser) blockStmt() (ast.Block, error) {
	_, err := p.eat(tokentype.LEFT_BRACE)

	if err != nil {
		return ast.Block{}, err
	}

	stmts := make([]ast.Node, 0)

	for !p.isThenEat(tokentype.RIGHT_BRACE) {
		stmt, err := p.stmtOrExpr()

		if err != nil {
			return ast.Block{}, err
		}

		stmts = append(stmts, stmt)
	}

	return ast.Block{Contents: stmts}, nil
}

func (p *parser) retStmt() (ast.Node, error) {
	var expr ast.Expr
	var err error

	if !p.is(tokentype.SEMICOLON) {
		expr, err = p.expression()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.eat(tokentype.SEMICOLON)

	if err != nil {
		return nil, err
	}

	return ast.ReturnStmt{Expr: expr}, nil
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

	return p.call()
}

func (p *parser) call() (ast.Expr, error) {
	maybeIdent, err := p.fundamental()

	if err != nil {
		return nil, err
	}

	// If the next token is not an open parenthesis we do not
	// have a call expression so return whatever we got from
	// `p.fundamental()`
	if p.current().Type != tokentype.LEFT_PAREN {
		return maybeIdent, nil
	}

	callee := ast.CallExpr{Callee: maybeIdent}

	// To allow chained call expression like `someFn()()`
	// we loop through so long as we have a left paren to
	// consume and nest the functions together.
	for p.isThenEat(tokentype.LEFT_PAREN) {
		args := []ast.Expr{}

		for !p.isThenEat(tokentype.RIGHT_PAREN) {
			arg, err := p.expression()

			if err != nil {
				return nil, err
			}

			args = append(args, arg)

			// In case there are multiple arguments, consume the next comma
			p.isThenEat(tokentype.COMMA)
		}

		callee.Arguments = args
		callee = ast.CallExpr{Callee: callee}
	}

	return callee.Callee, nil
}

func (p *parser) function() (ast.Expr, error) {
	// Get formal parameter list
	_, err := p.eat(tokentype.LEFT_PAREN)

	if err != nil {
		return nil, err
	}

	var idents []ast.Identifier

	if !p.is(tokentype.RIGHT_PAREN) {
		idents, err = p.varDeclarationNames()

		if err != nil {
			return nil, err
		}
	}

	_, err = p.eat(tokentype.RIGHT_PAREN)

	if err != nil {
		return nil, err
	}

	var body ast.Block

	// Get function body
	if p.isThenEat(tokentype.MINUS_GREAT) {
		expr, err := p.expression()

		if err != nil {
			return nil, err
		}

		body.Contents = []ast.Node{
			ast.ReturnStmt{Expr: expr},
		}
	} else {
		body, err = p.blockStmt()

		if err != nil {
			return nil, err
		}
	}

	return ast.FuncExpr{Params: idents, Body: body}, nil
}

func (p *parser) tuple() (ast.Expr, error) {
	if p.isThenEat(tokentype.RIGHT_BRACKET) {
		return ast.TupleLiteralExpr{}, nil
	}

	items, err := p.commaExpressions()

	if err != nil {
		return nil, err
	}

	_, err = p.eat(tokentype.RIGHT_BRACKET)

	if err != nil {
		return nil, err
	}

	return ast.TupleLiteralExpr{Contents: items}, nil
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

	if p.isThenEat(tokentype.FN) {
		return p.function()
	}

	if p.isThenEat(tokentype.LEFT_BRACKET) {
		return p.tuple()
	}

	t := p.tokens[p.i]
	return nil, errors.ParseError{Msg: fmt.Sprintf("Malformed expression at %d: %d", t.Position.Row, t.Position.Col)}
}

func New(ts []tokens.Token) *parser {
	return &parser{tokens: ts}
}
