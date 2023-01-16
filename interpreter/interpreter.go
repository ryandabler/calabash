package interpreter

import (
	"calabash/ast"
	"calabash/errors"
	"calabash/internal/tokentype"
	"calabash/internal/visitor"
	"fmt"
	"math"
	"strconv"
)

type interpreter struct{}

func (i *interpreter) Eval(e ast.Expr) (interface{}, error) {
	return visitor.Accept[interface{}](e, i)
}

func (i *interpreter) VisitBinaryExpr(e ast.BinaryExpr) (interface{}, error) {
	l, err := i.Eval(e.Left)

	if err != nil {
		return nil, err
	}

	r, err := i.Eval(e.Right)

	if err != nil {
		return nil, err
	}

	op := e.Operator.Type

	switch op {
	case tokentype.PLUS:
		{
			ln, okl := l.(vnumber)
			rn, okr := r.(vnumber)

			if okl && okr {
				val := vnumber{
					value: ln.value + rn.value,
				}

				return val, nil
			}

			ls, okl := l.(vstring)
			rs, okr := r.(vstring)

			if okl && okr {
				val := vstring{
					value: ls.value + rs.value,
				}

				return val, nil
			}

			return nil, errors.RuntimeError{Msg: "The types for binary '+' are not the same"}
		}

	case tokentype.MINUS:
		{
			ln, okl := l.(vnumber)
			rn, okr := r.(vnumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '-' can only be performed on numbers"}
			}

			val := vnumber{
				value: ln.value - rn.value,
			}

			return val, nil
		}

	case tokentype.ASTERISK:
		{
			ln, okl := l.(vnumber)
			rn, okr := r.(vnumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '*' can only be performed on numbers"}
			}

			val := vnumber{
				value: ln.value * rn.value,
			}

			return val, nil
		}

	case tokentype.SLASH:
		{
			ln, okl := l.(vnumber)
			rn, okr := r.(vnumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '/' can only be performed on numbers"}
			}

			val := vnumber{
				value: ln.value / rn.value,
			}

			return val, nil
		}

	case tokentype.ASTERISK_ASTERISK:
		{
			ln, okl := l.(vnumber)
			rn, okr := r.(vnumber)

			if !okl || !okr {
				return nil, errors.RuntimeError{Msg: "Binary '**' can only be performed on numbers"}
			}

			val := vnumber{
				value: math.Pow(ln.value, rn.value),
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

	v := vnumber{
		value: n,
	}

	return v, nil
}

func (i *interpreter) VisitStrLitExpr(e ast.StringLiteralExpr) (interface{}, error) {
	rs := []rune(e.Value.Lexeme)
	l := len(rs)
	str := vstring{value: string(rs[1 : l-1])}

	return str, nil
}

func (i *interpreter) VisitGroupingExpr(e ast.GroupingExpr) (interface{}, error) {
	return i.Eval(e.Expr)
}

func (i *interpreter) VisitUnaryExpr(e ast.UnaryExpr) (interface{}, error) {
	expr, err := i.Eval(e.Expr)

	if err != nil {
		return nil, err
	}

	op := e.Operator.Type

	switch op {
	case tokentype.MINUS:
		{
			if val, ok := expr.(vnumber); ok {
				val.value *= -1
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
