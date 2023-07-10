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

const value vtype = iota

type Value interface {
	v() vtype
	Hash() string
}

type Number struct {
	Value float64
}

func (v *Number) v() vtype {
	return value
}

func (v *Number) Hash() string {
	return fmt.Sprintf("n:%v", v.Value)
}

type String struct {
	Value string
}

func (v *String) v() vtype {
	return value
}

func (v *String) Hash() string {
	return fmt.Sprintf("s:%q", v.Value)
}

type Bottom struct{}

func (v *Bottom) v() vtype {
	return value
}

func (v *Bottom) Hash() string {
	return "btm"
}

type Boolean struct {
	Value bool
}

func (v *Boolean) v() vtype {
	return value
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
	return value
}

func (v *Function) Hash() string {
	if v.hash == "" {
		// Because functions are always unique, to populate the unexported
		// hash we need to manually construct functions in this package
		// with the UUID supplied
		v.hash = uuid.V4()
	}

	return v.hash
}

func (v *Function) Apply(vs []Value) *Function {
	return &Function{
		Params: v.Params,
		Body:   v.Body,
		Apps:   append(v.Apps, vs...),
	}
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

type Tuple struct {
	Items []Value
	hash  string
}

func (v *Tuple) v() vtype {
	return value
}

func (v *Tuple) Hash() string {
	if v.hash == "" {
		v.hash = fmt.Sprintf("tpl:%s", slice.Fold(v.Items, "", func(i Value, acc string) string {
			return acc + "," + i.Hash()
		}))
	}

	return v.hash
}

type Proto struct {
	Methods map[string]*Function
	Keys    []string
	hash    string
}

func (v *Proto) v() vtype {
	return value
}

func (v *Proto) Hash() string {
	if v.hash == "" {
		v.hash = slice.Fold(v.Keys, "prt:", func(k string, acc string) string {
			return acc + k + "->" + v.Methods[k].Hash()
		})
	}

	return v.hash
}
