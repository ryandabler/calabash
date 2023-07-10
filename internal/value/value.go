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
	proto
)

type Value interface {
	v() vtype
	Hash() string
}

type Number struct {
	Value float64
}

func (v *Number) v() vtype {
	return num
}

func (v *Number) Hash() string {
	return fmt.Sprintf("n:%v", v)
}

type String struct {
	Value string
}

func (v *String) v() vtype {
	return str
}

func (v *String) Hash() string {
	return fmt.Sprintf("s:%q", v.Value)
}

type Bottom struct{}

func (v *Bottom) v() vtype {
	return bottom
}

func (v *Bottom) Hash() string {
	return "btm"
}

type Boolean struct {
	Value bool
}

func (v *Boolean) v() vtype {
	return boolean
}

func (v *Boolean) Hash() string {
	return fmt.Sprintf("b:%t", v.Value)
}

type Function struct {
	Params []ast.Identifier
	Body   ast.Block
	Apps   []Value
	hash   string
}

func (v *Function) v() vtype {
	return fn
}

func (v *Function) Hash() string {
	return v.hash
}

func (v *Function) Apply(vs []Value) *Function {
	f := NewFunction()
	f.Params = v.Params
	f.Body = v.Body
	f.Apps = append(v.Apps, vs...)

	return f
}

func (v *Function) Arity() int {
	return len(v.Params) - len(v.Apps)
}

func (v *Function) Call(e Evaluator) (interface{}, error) {
	rVal, err := e.Eval(v.Body.Contents)

	if err != nil {
		return nil, err
	}

	if rVal == nil {
		return Bottom{}, nil
	}

	return rVal, nil
}

// Because functions are always unique, to populate the unexported
// hash we need to manually construct functions in this package
// with the UUID supplied
func NewFunction() *Function {
	return &Function{hash: uuid.V4()}
}

type Tuple struct {
	Items []Value
	hash  string
}

func (v *Tuple) v() vtype {
	return tuple
}

func (v *Tuple) Hash() string {
	return fmt.Sprintf("tpl:%s", v.hash)
}

func NewTuple(items []Value) *Tuple {
	return &Tuple{
		Items: items,
		hash: slice.Fold(items, "", func(i Value, acc string) string {
			return acc + "," + i.Hash()
		}),
	}
}

type Proto struct {
	Methods map[string]*Function
	keys    []string
	hash    string
}

func (v *Proto) v() vtype {
	return proto
}

func (v *Proto) Hash() string {
	return v.hash
}

func NewProto(ks []string, ms map[string]*Function) *Proto {
	hash := "proto:"

	for _, k := range ks {
		hash += fmt.Sprintf("%s->%s,", k, ms[k].Hash())
	}

	return &Proto{
		Methods: ms,
		keys:    ks,
		hash:    hash,
	}
}
