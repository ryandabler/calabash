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
}

func Accept[T any](e ast.Expr, v visitor[T]) (T, error) {
	var empty T

	if n, ok := e.(ast.BinaryExpr); ok {
		return v.VisitBinaryExpr(n)
	}

	if n, ok := e.(ast.UnaryExpr); ok {
		return v.VisitUnaryExpr(n)
	}

	if n, ok := e.(ast.GroupingExpr); ok {
		return v.VisitGroupingExpr(n)
	}

	if n, ok := e.(ast.NumericLiteralExpr); ok {
		return v.VisitNumLitExpr(n)
	}

	if n, ok := e.(ast.StringLiteralExpr); ok {
		return v.VisitStrLitExpr(n)
	}

	return empty, errors.New("Supplied node did not match any node type")
}
