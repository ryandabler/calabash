package ast

import "calabash/lexer/tokens"

type nodetype int

const (
	binary_expr nodetype = iota
	unary_expr
	numeric_lit_expr
	string_lit_expr
	bottom_lit_expr
	boolean_lit_expr
	identifer_expr
	grouping_expr
	var_decl_stmt
	ident
	assign_stmt
	if_stmt
	block
)

type Node interface {
	n() nodetype
}

type Expr interface {
	e() nodetype
	n() nodetype
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator tokens.Token
}

func (e BinaryExpr) e() nodetype {
	return binary_expr
}

func (e BinaryExpr) n() nodetype {
	return binary_expr
}

type UnaryExpr struct {
	Operator tokens.Token
	Expr     Expr
}

func (e UnaryExpr) e() nodetype {
	return unary_expr
}

func (e UnaryExpr) n() nodetype {
	return unary_expr
}

type NumericLiteralExpr struct {
	Value tokens.Token
}

func (e NumericLiteralExpr) e() nodetype {
	return numeric_lit_expr
}

func (e NumericLiteralExpr) n() nodetype {
	return numeric_lit_expr
}

type StringLiteralExpr struct {
	Value tokens.Token
}

func (e StringLiteralExpr) e() nodetype {
	return string_lit_expr
}

func (e StringLiteralExpr) n() nodetype {
	return string_lit_expr
}

type BottomLiteralExpr struct {
	Token tokens.Token
}

func (e BottomLiteralExpr) e() nodetype {
	return bottom_lit_expr
}

func (e BottomLiteralExpr) n() nodetype {
	return bottom_lit_expr
}

type BooleanLiteralExpr struct {
	Value tokens.Token
}

func (e BooleanLiteralExpr) e() nodetype {
	return boolean_lit_expr
}

func (e BooleanLiteralExpr) n() nodetype {
	return boolean_lit_expr
}

type IdentifierExpr struct {
	Name tokens.Token
}

func (e IdentifierExpr) e() nodetype {
	return identifer_expr
}

func (e IdentifierExpr) n() nodetype {
	return identifer_expr
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) e() nodetype {
	return grouping_expr
}

func (e GroupingExpr) n() nodetype {
	return grouping_expr
}

type VarDeclStmt struct {
	Names  []Identifier
	Values []Expr
}

func (s VarDeclStmt) n() nodetype {
	return var_decl_stmt
}

type Identifier struct {
	Name tokens.Token
	Mut  bool
}

func (s Identifier) n() nodetype {
	return ident
}

type AssignmentStmt struct {
	Names  []tokens.Token
	Values []Expr
}

func (s AssignmentStmt) n() nodetype {
	return assign_stmt
}

type IfStmt struct {
	Decls     VarDeclStmt
	Condition Expr
	Then      Node
	Else      Node
}

func (s IfStmt) n() nodetype {
	return if_stmt
}

type Block struct {
	Contents []Node
}

func (s Block) n() nodetype {
	return block
}
