package staticanalyzer

import (
	"calabash/ast"
	"calabash/errors"
	"calabash/internal/environment"
	"calabash/internal/visitor"
	"fmt"
)

type analyzer struct {
	env *environment.Environment[identRecord]
}

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

func (a *analyzer) VisitBooleanLitExpr(e ast.BooleanLiteralExpr) (interface{}, error) {
	return nil, nil
}

func (a *analyzer) VisitIdentifierExpr(e ast.IdentifierExpr) (interface{}, error) {
	if !a.env.Has(e.Name.Lexeme) {
		return nil, errors.StaticError{Msg: "Cannot reference an undeclared identifier."}
	}

	return nil, nil
}

func (a *analyzer) VisitVarDeclStmt(s ast.VarDeclStmt) (interface{}, error) {
	if len(s.Names) != len(s.Values) && len(s.Values) > 0 {
		return nil, errors.StaticError{Msg: "If any variable is initialized, they all must be."}
	}

	for _, n := range s.Names {
		if a.env.HasDirectly(n.Name.Lexeme) {
			return nil, errors.StaticError{Msg: fmt.Sprintf("Cannot redeclare variable %q", n.Name.Lexeme)}
		}

		a.env.Add(n.Name.Lexeme, identRecord{mut: n.Mut})
	}

	for _, v := range s.Values {
		err := a.analyzeNode(v)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (a *analyzer) VisitAssignStmt(s ast.AssignmentStmt) (interface{}, error) {
	if len(s.Names) != len(s.Values) {
		msg := fmt.Sprintf("Expected to have %d expressions--received %d", len(s.Names), len(s.Values))
		return nil, errors.StaticError{Msg: msg}
	}

	for _, n := range s.Names {
		if !a.env.Has(n.Lexeme) {
			return nil, errors.StaticError{Msg: fmt.Sprintf("Cannot assign to undeclared variable %q", n.Lexeme)}
		}

		if !a.env.Get(n.Lexeme).mut {
			return nil, errors.StaticError{Msg: fmt.Sprintf("Cannot re-assign immutable variable %q", n.Lexeme)}
		}
	}

	for _, v := range s.Values {
		err := a.analyzeNode(v)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (a *analyzer) VisitIfStmt(s ast.IfStmt) (interface{}, error) {
	// Set new environment for the entire level of the if-then-else blocks
	e := environment.New(a.env)
	a.env = e

	// Analyze any declarations
	if len(s.Decls.Names) > 0 {
		_, err := a.VisitVarDeclStmt(s.Decls)

		if err != nil {
			return nil, err
		}
	}

	// Analyze the `if`` condition
	err := a.analyzeNode(s.Condition)

	if err != nil {
		return nil, err
	}

	// Analyze the `then` block
	err = a.analyzeNode(s.Then)

	if err != nil {
		return nil, err
	}

	// If there is an `else` block, analyze it
	if s.Else != nil {
		err = a.analyzeNode(s.Else)

		if err != nil {
			return nil, err
		}
	}

	a.env = e.Parent

	return nil, nil
}

func (a *analyzer) VisitBlock(s ast.Block) (interface{}, error) {
	e := environment.New(a.env)
	a.env = e

	for _, n := range s.Contents {
		err := a.analyzeNode(n)

		if err != nil {
			return nil, err
		}
	}

	a.env = e.Parent

	return nil, nil
}

func New() *analyzer {
	return &analyzer{
		env: environment.New[identRecord](nil),
	}
}
