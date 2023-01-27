package value

type vtype int

const (
	num vtype = iota
	str
	bottom
)

type Value interface {
	v() vtype
}

type VNumber struct {
	Value float64
}

func (v VNumber) v() vtype {
	return num
}

type VString struct {
	Value string
}

func (v VString) v() vtype {
	return str
}

type VBottom struct{}

func (v VBottom) v() vtype {
	return bottom
}
