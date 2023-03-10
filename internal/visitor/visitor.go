package visitor

import (
	"calabash/ast"
	"errors"
)

type visitor[T any] interface {
	VisitBinaryExpr(e ast.BinaryExpr) (T, error)
	VisitUnaryExpr(e ast.UnaryExpr) (T, error)
	VisitGroupingExpr(e ast.GroupingExpr) (T, error)
	VisitNumLitExpr(e ast.NumericLiteralExpr) (T, error)
	VisitStrLitExpr(e ast.StringLiteralExpr) (T, error)
	VisitBottomLitExpr(e ast.BottomLiteralExpr) (T, error)
	VisitBooleanLitExpr(e ast.BooleanLiteralExpr) (T, error)
	VisitIdentifierExpr(e ast.IdentifierExpr) (T, error)
	VisitFuncExpr(e ast.FuncExpr) (T, error)
	VisitVarDeclStmt(s ast.VarDeclStmt) (T, error)
	VisitAssignStmt(s ast.AssignmentStmt) (T, error)
	VisitIfStmt(s ast.IfStmt) (T, error)
	VisitBlock(s ast.Block) (T, error)
	VisitRetStmt(s ast.ReturnStmt) (T, error)
}

func Accept[T any](n ast.Node, v visitor[T]) (T, error) {
	var empty T

	if n, ok := n.(ast.BinaryExpr); ok {
		return v.VisitBinaryExpr(n)
	}

	if n, ok := n.(ast.UnaryExpr); ok {
		return v.VisitUnaryExpr(n)
	}

	if n, ok := n.(ast.GroupingExpr); ok {
		return v.VisitGroupingExpr(n)
	}

	if n, ok := n.(ast.NumericLiteralExpr); ok {
		return v.VisitNumLitExpr(n)
	}

	if n, ok := n.(ast.StringLiteralExpr); ok {
		return v.VisitStrLitExpr(n)
	}

	if n, ok := n.(ast.BottomLiteralExpr); ok {
		return v.VisitBottomLitExpr(n)
	}

	if n, ok := n.(ast.BooleanLiteralExpr); ok {
		return v.VisitBooleanLitExpr(n)
	}

	if n, ok := n.(ast.IdentifierExpr); ok {
		return v.VisitIdentifierExpr(n)
	}

	if n, ok := n.(ast.FuncExpr); ok {
		return v.VisitFuncExpr(n)
	}

	if n, ok := n.(ast.VarDeclStmt); ok {
		return v.VisitVarDeclStmt(n)
	}

	if n, ok := n.(ast.AssignmentStmt); ok {
		return v.VisitAssignStmt(n)
	}

	if n, ok := n.(ast.IfStmt); ok {
		return v.VisitIfStmt(n)
	}

	if n, ok := n.(ast.Block); ok {
		return v.VisitBlock(n)
	}

	if n, ok := n.(ast.ReturnStmt); ok {
		return v.VisitRetStmt(n)
	}

	return empty, errors.New("Supplied node did not match any node type")
}
