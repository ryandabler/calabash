package value

import (
	"calabash/ast"
	"calabash/internal/environment"
)

type vtype int

const value vtype = iota

type Evaluator interface {
	Eval([]ast.Node) (interface{}, error)
	Dump() struct {
		Env *environment.Environment[Value]
	}
}

type Value interface {
	v() vtype
	Hash() string
}

type Caller interface {
	Apply([]Value) Caller
	Args() []Value
	Params() []ast.Identifier
	Arity() int
	Call(Evaluator) (interface{}, error)
	Hash() string
}
