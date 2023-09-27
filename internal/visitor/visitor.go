package visitor

import (
	"calabash/ast"
	"errors"
)

type evisitor[T any] interface {
	VisitBinaryExpr(e ast.BinaryExpr) (T, error)
	VisitUnaryExpr(e ast.UnaryExpr) (T, error)
	VisitGroupingExpr(e ast.GroupingExpr) (T, error)
	VisitNumLitExpr(e ast.NumericLiteralExpr) (T, error)
	VisitStrLitExpr(e ast.StringLiteralExpr) (T, error)
	VisitBottomLitExpr(e ast.BottomLiteralExpr) (T, error)
	VisitBooleanLitExpr(e ast.BooleanLiteralExpr) (T, error)
	VisitTupleLitExpr(e ast.TupleLiteralExpr) (T, error)
	VisitIdentifierExpr(e ast.IdentifierExpr) (T, error)
	VisitFuncExpr(e ast.FuncExpr) (T, error)
	VisitCallExpr(e ast.CallExpr) (T, error)
	VisitMeExpr(e ast.MeExpr) (T, error)
	VisitProtoExpr(e ast.ProtoExpr) (T, error)
	VisitGetExpr(e ast.GetExpr) (T, error)
}

type svisitor[T any] interface {
	VisitVarDeclStmt(s ast.VarDeclStmt) (T, error)
	VisitAssignStmt(s ast.AssignmentStmt) (T, error)
	VisitIfStmt(s ast.IfStmt) (T, error)
	VisitBlock(s ast.Block) (T, error)
	VisitRetStmt(s ast.ReturnStmt) (T, error)
	VisitWhileStmt(s ast.WhileStmt) (T, error)
	VisitContStmt(s ast.ContinueStmt) (T, error)
	VisitBrkStmt(s ast.BreakStmt) (T, error)
}

type visitor[T any] interface {
	evisitor[T]
	svisitor[T]
}

func AcceptExpr[T any](e ast.Expr, v evisitor[T]) (T, error) {
	var empty T

	switch e.(type) {
	case ast.BinaryExpr:
		e := e.(ast.BinaryExpr)

		return v.VisitBinaryExpr(e)

	case ast.UnaryExpr:
		e := e.(ast.UnaryExpr)

		return v.VisitUnaryExpr(e)

	case ast.GroupingExpr:
		e := e.(ast.GroupingExpr)

		return v.VisitGroupingExpr(e)

	case ast.NumericLiteralExpr:
		e := e.(ast.NumericLiteralExpr)

		return v.VisitNumLitExpr(e)

	case ast.StringLiteralExpr:
		e := e.(ast.StringLiteralExpr)

		return v.VisitStrLitExpr(e)

	case ast.BottomLiteralExpr:
		e := e.(ast.BottomLiteralExpr)

		return v.VisitBottomLitExpr(e)

	case ast.BooleanLiteralExpr:
		e := e.(ast.BooleanLiteralExpr)

		return v.VisitBooleanLitExpr(e)

	case ast.TupleLiteralExpr:
		e := e.(ast.TupleLiteralExpr)

		return v.VisitTupleLitExpr(e)

	case ast.IdentifierExpr:
		e := e.(ast.IdentifierExpr)

		return v.VisitIdentifierExpr(e)

	case ast.FuncExpr:
		e := e.(ast.FuncExpr)

		return v.VisitFuncExpr(e)

	case ast.CallExpr:
		e, ok := e.(ast.CallExpr)

		if !ok {
			return empty, errors.New("Count not convert to call expr")
		}

		return v.VisitCallExpr(e)

	case ast.MeExpr:
		e := e.(ast.MeExpr)

		return v.VisitMeExpr(e)

	case ast.ProtoExpr:
		e := e.(ast.ProtoExpr)

		return v.VisitProtoExpr(e)

	case ast.GetExpr:
		e := e.(ast.GetExpr)

		return v.VisitGetExpr(e)
	}

	return empty, errors.New("Unexpected expression")
}

func Accept[T any](n ast.Node, v visitor[T]) (T, error) {
	var empty T

	e, ok := n.(ast.Expr)

	if ok {
		return AcceptExpr[T](e, v)
	}

	switch n.(type) {
	case ast.VarDeclStmt:
		s := n.(ast.VarDeclStmt)

		return v.VisitVarDeclStmt(s)

	case ast.AssignmentStmt:
		s := n.(ast.AssignmentStmt)

		return v.VisitAssignStmt(s)

	case ast.IfStmt:
		s := n.(ast.IfStmt)

		return v.VisitIfStmt(s)

	case ast.Block:
		s := n.(ast.Block)

		return v.VisitBlock(s)

	case ast.ReturnStmt:
		s := n.(ast.ReturnStmt)

		return v.VisitRetStmt(s)

	case ast.WhileStmt:
		s := n.(ast.WhileStmt)

		return v.VisitWhileStmt(s)

	case ast.ContinueStmt:
		s := n.(ast.ContinueStmt)

		return v.VisitContStmt(s)

	case ast.BreakStmt:
		s := n.(ast.BreakStmt)

		return v.VisitBrkStmt(s)
	}

	return empty, errors.New("Supplied node did not match any node type")
}
