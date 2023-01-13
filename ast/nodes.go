package ast

import "calabash/lexer/tokens"

type nodetype int

const (
	binary_expr nodetype = iota
	unary_expr
	numeric_lit_expr
	string_lit_expr
	grouping_expr
)

type Expr interface {
	e() nodetype
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator tokens.Token
}

func (e BinaryExpr) e() nodetype {
	return binary_expr
}

type UnaryExpr struct {
	Operator tokens.Token
	Expr     Expr
}

func (e UnaryExpr) e() nodetype {
	return unary_expr
}

type NumericLiteralExpr struct {
	Value tokens.Token
}

func (e NumericLiteralExpr) e() nodetype {
	return numeric_lit_expr
}

type StringLiteralExpr struct {
	Value tokens.Token
}

func (e StringLiteralExpr) e() nodetype {
	return string_lit_expr
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) e() nodetype {
	return grouping_expr
}
