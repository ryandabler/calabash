package interpreter

import (
	"calabash/ast"
	"calabash/errors"
	"calabash/internal/tokentype"
	"calabash/internal/value"
	"calabash/internal/visitor"
	"fmt"
	"math"
	"strconv"
)

type interpreter struct{}

func (i *interpreter) Eval(ns []ast.Node) (interface{}, error) {
	var v interface{}
	var err error

	for _, n := range ns {
		v, err = i.evalNode(n)

		if err != nil {
			return nil, err
		}
	}

	return v, err
}

func (i *interpreter) evalNode(n ast.Node) (interface{}, error) {
	return visitor.Accept[interface{}](n, i)
}

func (i *interpreter) VisitBinaryExpr(e ast.BinaryExpr) (interface{}, error) {
	l, err := i.evalNode(e.Left)

	if err != nil {
		return nil, err
	}

	r, err := i.evalNode(e.Right)

	if err != nil {
		return nil, err
	}

	op := e.Operator.Type

	switch op {
	case tokentype.PLUS:
		{
			ln, okl := l.(value.VNumber)
			rn, okr := r.(value.VNumber)

			if okl && okr {
				val := value.VNumber{
					Value: ln.Value + rn.Value,
				}

				return val, nil
			}

			ls, okl := l.(value.VString)
			rs, okr := r.(value.VString)

			if okl && okr {
				val := value.VString{
					Value: ls.Value + rs.Value,
				}

				return val, nil
			}

			return nil, errors.RuntimeError{Msg: "The types for binary '+' are not the same"}
		}

	case tokentype.MINUS:
		{
			ln, okl := l.(value.VNumber)
			rn, okr := r.(value.VNumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '-' can only be performed on numbers"}
			}

			val := value.VNumber{
				Value: ln.Value - rn.Value,
			}

			return val, nil
		}

	case tokentype.ASTERISK:
		{
			ln, okl := l.(value.VNumber)
			rn, okr := r.(value.VNumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '*' can only be performed on numbers"}
			}

			val := value.VNumber{
				Value: ln.Value * rn.Value,
			}

			return val, nil
		}

	case tokentype.SLASH:
		{
			ln, okl := l.(value.VNumber)
			rn, okr := r.(value.VNumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '/' can only be performed on numbers"}
			}

			val := value.VNumber{
				Value: ln.Value / rn.Value,
			}

			return val, nil
		}

	case tokentype.ASTERISK_ASTERISK:
		{
			ln, okl := l.(value.VNumber)
			rn, okr := r.(value.VNumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '**' can only be performed on numbers"}
			}

			val := value.VNumber{
				Value: math.Pow(ln.Value, rn.Value),
			}

			return val, nil
		}
	}

	return nil, errors.RuntimeError{Msg: fmt.Sprintf("The only supported binary operators are '+': received %q", e.Operator.Lexeme)}
}

func (i *interpreter) VisitNumLitExpr(e ast.NumericLiteralExpr) (interface{}, error) {
	n, err := strconv.ParseFloat(e.Value.Lexeme, 64)

	if err != nil {
		return nil, err
	}

	v := value.VNumber{
		Value: n,
	}

	return v, nil
}

func (i *interpreter) VisitStrLitExpr(e ast.StringLiteralExpr) (interface{}, error) {
	rs := []rune(e.Value.Lexeme)
	l := len(rs)
	str := value.VString{Value: string(rs[1 : l-1])}

	return str, nil
}

func (i *interpreter) VisitGroupingExpr(e ast.GroupingExpr) (interface{}, error) {
	return i.evalNode(e.Expr)
}

func (i *interpreter) VisitUnaryExpr(e ast.UnaryExpr) (interface{}, error) {
	expr, err := i.evalNode(e.Expr)

	if err != nil {
		return nil, err
	}

	op := e.Operator.Type

	switch op {
	case tokentype.MINUS:
		{
			if val, ok := expr.(value.VNumber); ok {
				val.Value *= -1
				return val, nil
			}

			return nil, errors.RuntimeError{Msg: "Can only use unary minus with numbers."}
		}
	}

	return nil, errors.RuntimeError{Msg: fmt.Sprintf("The only supported unary operators are '-': got %q", e.Operator.Lexeme)}
}

func New() *interpreter {
	return &interpreter{}
}
