package value

import (
	"calabash/ast"
	"calabash/internal/slice"
	"calabash/internal/uuid"
	"fmt"
)

type Evaluator interface {
	Eval([]ast.Node) (interface{}, error)
}

type vtype int

const (
	num vtype = iota
	str
	boolean
	bottom
	fn
	tuple
)

type Value interface {
	v() vtype
	Hash() string
}

type VNumber struct {
	Value float64
}

func (v VNumber) v() vtype {
	return num
}

func (v VNumber) Hash() string {
	return fmt.Sprintf("n:%v", v)
}

type VString struct {
	Value string
}

func (v VString) v() vtype {
	return str
}

func (v VString) Hash() string {
	return fmt.Sprintf("s:%q", v.Value)
}

type VBottom struct{}

func (v VBottom) v() vtype {
	return bottom
}

func (v VBottom) Hash() string {
	return "btm"
}

type VBoolean struct {
	Value bool
}

func (v VBoolean) v() vtype {
	return boolean
}

func (v VBoolean) Hash() string {
	return fmt.Sprintf("b:%t", v.Value)
}

type VFunction struct {
	Params []ast.Identifier
	Body   ast.Block
	Apps   []Value
	hash   string
}

func (v VFunction) v() vtype {
	return fn
}

func (v VFunction) Hash() string {
	return v.hash
}

func (v VFunction) Apply(vs []Value) VFunction {
	f := NewFunction()
	f.Params = v.Params
	f.Body = v.Body
	f.Apps = append(v.Apps, vs...)

	return f
}

func (v VFunction) Arity() int {
	return len(v.Params) - len(v.Apps)
}

func (v VFunction) Call(e Evaluator) (interface{}, error) {
	rVal, err := e.Eval(v.Body.Contents)

	if err != nil {
		return nil, err
	}

	if rVal == nil {
		return VBottom{}, nil
	}

	return rVal, nil
}

// Because functions are always unique, to populate the unexported
// hash we need to manually construct functions in this package
// with the UUID supplied
func NewFunction() VFunction {
	return VFunction{hash: uuid.V4()}
}

type VTuple struct {
	Items []Value
	hash  string
}

func (v VTuple) v() vtype {
	return tuple
}

func (v VTuple) Hash() string {
	return fmt.Sprintf("tpl:%s", v.hash)
}

func NewTuple(items []Value) VTuple {
	return VTuple{
		Items: items,
		hash: slice.Fold(items, "", func(i Value, acc string) string {
			return acc + "," + i.Hash()
		}),
	}
}
