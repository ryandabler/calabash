package ast

import "calabash/lexer/tokens"

type nodetype int

const nt nodetype = iota

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
	return nt
}

func (e BinaryExpr) n() nodetype {
	return nt
}

type UnaryExpr struct {
	Operator tokens.Token
	Expr     Expr
}

func (e UnaryExpr) e() nodetype {
	return nt
}

func (e UnaryExpr) n() nodetype {
	return nt
}

type NumericLiteralExpr struct {
	Value tokens.Token
}

func (e NumericLiteralExpr) e() nodetype {
	return nt
}

func (e NumericLiteralExpr) n() nodetype {
	return nt
}

type StringLiteralExpr struct {
	Value tokens.Token
}

func (e StringLiteralExpr) e() nodetype {
	return nt
}

func (e StringLiteralExpr) n() nodetype {
	return nt
}

type BottomLiteralExpr struct {
	Token tokens.Token
}

func (e BottomLiteralExpr) e() nodetype {
	return nt
}

func (e BottomLiteralExpr) n() nodetype {
	return nt
}

type BooleanLiteralExpr struct {
	Value tokens.Token
}

func (e BooleanLiteralExpr) e() nodetype {
	return nt
}

func (e BooleanLiteralExpr) n() nodetype {
	return nt
}

type TupleLiteralExpr struct {
	Contents []Expr
}

func (e TupleLiteralExpr) e() nodetype {
	return nt
}

func (e TupleLiteralExpr) n() nodetype {
	return nt
}

type IdentifierExpr struct {
	Name tokens.Token
}

func (e IdentifierExpr) e() nodetype {
	return nt
}

func (e IdentifierExpr) n() nodetype {
	return nt
}

type GroupingExpr struct {
	Expr Expr
}

func (e GroupingExpr) e() nodetype {
	return nt
}

func (e GroupingExpr) n() nodetype {
	return nt
}

type FuncExpr struct {
	Params []Identifier
	Body   Block
}

func (e FuncExpr) e() nodetype {
	return nt
}

func (e FuncExpr) n() nodetype {
	return nt
}

type CallExpr struct {
	Callee    Expr
	Arguments []Expr
}

func (e CallExpr) e() nodetype {
	return nt
}

func (e CallExpr) n() nodetype {
	return nt
}

type MeExpr struct {
	Token tokens.Token
}

func (e MeExpr) e() nodetype {
	return nt
}

func (e MeExpr) n() nodetype {
	return nt
}

type ProtoExpr struct {
	MethodSet []ProtoMethod
}

func (e ProtoExpr) e() nodetype {
	return nt
}

func (e ProtoExpr) n() nodetype {
	return nt
}

type GetExpr struct {
	Gettee Expr
	Field  Expr
}

func (e GetExpr) e() nodetype {
	return nt
}

func (e GetExpr) n() nodetype {
	return nt
}

type VarDeclStmt struct {
	Names  []Identifier
	Values []Expr
}

func (s VarDeclStmt) n() nodetype {
	return nt
}

type Identifier struct {
	Name tokens.Token
	Mut  bool
}

func (s Identifier) n() nodetype {
	return nt
}

type AssignmentStmt struct {
	Names  []tokens.Token
	Values []Expr
}

func (s AssignmentStmt) n() nodetype {
	return nt
}

type IfStmt struct {
	Decls     VarDeclStmt
	Condition Expr
	Then      Node
	Else      Node
}

func (s IfStmt) n() nodetype {
	return nt
}

type Block struct {
	Contents []Node
}

func (s Block) n() nodetype {
	return nt
}

type ReturnStmt struct {
	Expr Expr
}

func (s ReturnStmt) n() nodetype {
	return nt
}

type WhileStmt struct {
	Decls     VarDeclStmt
	Condition Expr
	Block     Node
}

func (s WhileStmt) n() nodetype {
	return nt
}

type ContinueStmt struct{}

func (s ContinueStmt) n() nodetype {
	return nt
}

type BreakStmt struct{}

func (s BreakStmt) n() nodetype {
	return nt
}
