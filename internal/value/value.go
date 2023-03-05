package value

import (
	"calabash/ast"
	"calabash/internal/uuid"
	"fmt"
)

type vtype int

const (
	num vtype = iota
	str
	boolean
	bottom
	fn
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
	hash   string
}

func (v VFunction) v() vtype {
	return fn
}

func (v *VFunction) Hash() string {
	if v.hash == "" {
		v.hash = uuid.V4()
	}

	return v.hash
}
