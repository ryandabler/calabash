package staticanalyzer

import (
	"calabash/ast"
	"calabash/internal/visitor"
)

type analyzer struct{}

func (a *analyzer) Analyze(ast []ast.Node) error {
	for _, n := range ast {
		err := a.analyzeNode(n)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *analyzer) analyzeNode(n ast.Node) error {
	_, err := visitor.Accept[interface{}](n, a)

	return err
}

func (a *analyzer) VisitBinaryExpr(e ast.BinaryExpr) (interface{}, error) {
	err := a.analyzeNode(e.Left)

	if err != nil {
		return nil, err
	}

	err = a.analyzeNode(e.Right)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *analyzer) VisitGroupingExpr(e ast.GroupingExpr) (interface{}, error) {
	return nil, a.analyzeNode(e.Expr)
}

func (a *analyzer) VisitNumLitExpr(e ast.NumericLiteralExpr) (interface{}, error) {
	return nil, nil
}

func (a *analyzer) VisitStrLitExpr(e ast.StringLiteralExpr) (interface{}, error) {
	return nil, nil
}

func (a *analyzer) VisitUnaryExpr(e ast.UnaryExpr) (interface{}, error) {
	return nil, a.analyzeNode(e.Expr)
}

func (a *analyzer) VisitBottomLitExpr(e ast.BottomLiteralExpr) (interface{}, error) {
	return nil, nil
}

func New() *analyzer {
	return &analyzer{}
}
