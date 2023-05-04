package interpreter

import (
	"calabash/ast"
	"calabash/errors"
	"calabash/internal/environment"
	"calabash/internal/slice"
	"calabash/internal/tokentype"
	"calabash/internal/value"
	"calabash/internal/visitor"
	errs "errors"
	"fmt"
	"math"
	"strconv"
)

type interpreter struct {
	env *environment.Environment[value.Value]
}

func (i *interpreter) Eval(ns []ast.Node) (interface{}, error) {
	var v interface{}
	var err error

	for _, n := range ns {
		v, err = i.evalNode(n)

		if errs.Is(err, errors.ReturnError{}) {
			return v, nil
		}

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

	if op == tokentype.EQUAL_EQUAL || op == tokentype.BANG_EQUAL {
		l, okl := l.(value.Value)
		r, okr := r.(value.Value)

		if !okl || !okr {
			return nil, errors.RuntimeError{Msg: "Can not check equality of non-values"}
		}

		b := l.Hash() == r.Hash()

		if op == tokentype.BANG_EQUAL {
			b = !b
		}

		return value.VBoolean{Value: b}, nil
	}

	if isBooleanOp(op) && areBools(l, r) {
		lb, _ := l.(value.VBoolean)
		rb, _ := r.(value.VBoolean)

		if op == tokentype.AMPERSAND_AMPERSAND {
			return value.VBoolean{Value: lb.Value && rb.Value}, nil
		}

		if op == tokentype.STROKE_STROKE {
			return value.VBoolean{Value: lb.Value || rb.Value}, nil
		}
	}

	if isBooleanOp(op) {
		return nil, errors.RuntimeError{Msg: "Both operands must be boolean values to use with boolean operators."}
	}

	// The '+' operator is overloaded for different data types. The left and right
	// sides must be of the same type but they could be many different types.
	if op == tokentype.PLUS {
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

	if isNumericOp(op) && areNumbers(l, r) {
		ln, _ := l.(value.VNumber)
		rn, _ := r.(value.VNumber)
		var val value.Value

		switch op {
		case tokentype.MINUS:
			val = value.VNumber{
				Value: ln.Value - rn.Value,
			}

		case tokentype.ASTERISK:
			val = value.VNumber{
				Value: ln.Value * rn.Value,
			}

		case tokentype.SLASH:
			val = value.VNumber{
				Value: ln.Value / rn.Value,
			}

		case tokentype.ASTERISK_ASTERISK:
			val = value.VNumber{
				Value: math.Pow(ln.Value, rn.Value),
			}

		case tokentype.GREAT:
			val = value.VBoolean{Value: ln.Value > rn.Value}

		case tokentype.GREAT_EQUAL:
			val = value.VBoolean{Value: ln.Value >= rn.Value}

		case tokentype.LESS:
			val = value.VBoolean{Value: ln.Value < rn.Value}

		case tokentype.LESS_EQUAL:
			val = value.VBoolean{Value: ln.Value <= rn.Value}
		}

		return val, nil
	}

	if isNumericOp(op) {
		return nil, errors.RuntimeError{Msg: fmt.Sprintf("Received a non-numeric value for numeric binary operator %q", e.Operator.Lexeme)}
	}

	return nil, errors.RuntimeError{Msg: fmt.Sprintf("Received unsupported binary operator %q", e.Operator.Lexeme)}
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

func (i *interpreter) VisitBottomLitExpr(e ast.BottomLiteralExpr) (interface{}, error) {
	return value.VBottom{}, nil
}

func (i *interpreter) VisitBooleanLitExpr(e ast.BooleanLiteralExpr) (interface{}, error) {
	return value.VBoolean{Value: e.Value.Type == tokentype.TRUE}, nil
}

func (i *interpreter) VisitTupleLitExpr(e ast.TupleLiteralExpr) (interface{}, error) {
	vs := make([]value.Value, len(e.Contents))

	for idx, c := range e.Contents {
		ifc, err := i.evalNode(c)

		if err != nil {
			return nil, err
		}

		v, ok := ifc.(value.Value)

		if !ok {
			return nil, errs.New("Did not have a value when building tuple literal's contents")
		}

		vs[idx] = v
	}

	return value.NewTuple(vs), nil
}

func (i *interpreter) VisitIdentifierExpr(e ast.IdentifierExpr) (interface{}, error) {
	return i.env.Get(e.Name.Lexeme), nil
}

func (i *interpreter) VisitFuncExpr(e ast.FuncExpr) (interface{}, error) {
	fn := value.NewFunction()
	fn.Body = e.Body
	fn.Params = e.Params

	return fn, nil
}

func (i *interpreter) VisitCallExpr(e ast.CallExpr) (interface{}, error) {
	// Evaluate callee and arguments in current scope before replacing it
	// with empty scope for the function body.
	callee, err := i.evalNode(e.Callee)

	if err != nil {
		return nil, err
	}

	vfunc, ok := callee.(value.VFunction)

	if !ok {
		return nil, errors.RuntimeError{Msg: "Attempting to call a non-functional value."}
	}

	vals, err := slice.Map(e.Arguments, func(e ast.Expr) (value.Value, error) {
		v, err := i.evalNode(e)

		if err != nil {
			return nil, err
		}

		val, ok := v.(value.Value)

		if !ok {
			return nil, errors.RuntimeError{Msg: "Argument " + fmt.Sprint(1) + " is not a value"}
		}

		return val, nil
	})

	if err != nil {
		return nil, err
	}

	// Check if function is being partially applied and return a new function if so
	if len(vals) < vfunc.Arity() {
		return vfunc.Apply(vals), nil
	}

	// Begin function call routines
	fBodyEnv := environment.New[value.Value](nil)
	args := append(vfunc.Apps, vals...)

	for idx, ident := range vfunc.Params {
		fBodyEnv.Add(ident.Name.Lexeme, args[idx])
	}

	// By default, functions are not closures so they only have access to their
	// own environment
	env := i.env
	i.env = fBodyEnv

	v, err := vfunc.Call(i)

	if err != nil {
		return nil, err
	}

	i.env = env

	return v, nil
}

func (i *interpreter) VisitVarDeclStmt(s ast.VarDeclStmt) (interface{}, error) {
	for idx, n := range s.Names {
		var val value.Value = value.VBottom{}
		var ok bool

		// If initial values have been specified for these declarations,
		// determine them before adding them to the environment
		if len(s.Values) > 0 {
			v, err := i.evalNode(s.Values[idx])

			if err != nil {
				return nil, err
			}

			val, ok = v.(value.Value)

			if !ok {
				return nil, errors.RuntimeError{Msg: "Did not receive a value in variable declaration initialization."}
			}
		}

		i.env.Add(n.Name.Lexeme, val)
	}

	return nil, nil
}

func (i *interpreter) VisitAssignStmt(s ast.AssignmentStmt) (interface{}, error) {
	for idx, n := range s.Names {
		v, err := i.evalNode(s.Values[idx])

		if err != nil {
			return nil, err
		}

		val, ok := v.(value.Value)

		if !ok {
			return nil, errors.RuntimeError{Msg: "Could not obtain a value for assignment."}
		}

		i.env.Set(n.Lexeme, val)
	}

	return nil, nil
}

func (i *interpreter) VisitIfStmt(s ast.IfStmt) (interface{}, error) {
	e := environment.New(i.env)
	i.env = e

	defer func() { i.env = i.env.Parent }()

	_, err := i.VisitVarDeclStmt(s.Decls)

	if err != nil {
		return nil, err
	}

	vCond, err := i.evalNode(s.Condition)

	if err != nil {
		return nil, err
	}

	cond, ok := vCond.(value.VBoolean)

	if !ok {
		return nil, errors.RuntimeError{Msg: "If condition must resolve to a boolean value."}
	}

	if cond.Value {
		_, err = i.evalNode(s.Then)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	if s.Else == nil {
		return nil, nil
	}

	_, err = i.evalNode(s.Else)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *interpreter) VisitBlock(s ast.Block) (interface{}, error) {
	e := environment.New(i.env)
	i.env = e

	defer func() { i.env = i.env.Parent }()

	for _, n := range s.Contents {
		_, err := i.evalNode(n)

		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *interpreter) VisitRetStmt(s ast.ReturnStmt) (interface{}, error) {
	v, err := visitor.Accept[interface{}](s.Expr, i)

	if err != nil {
		return nil, err
	}

	return v, errors.ReturnError{}
}

func New() *interpreter {
	return &interpreter{
		env: environment.New[value.Value](nil),
	}
}

type IntpState struct {
	Env *environment.Environment[value.Value]
}

func (i *interpreter) Dump() IntpState {
	return IntpState{Env: i.env}
}
