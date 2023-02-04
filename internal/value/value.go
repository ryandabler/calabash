package value

import "fmt"

type vtype int

const (
	num vtype = iota
	str
	boolean
	bottom
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
