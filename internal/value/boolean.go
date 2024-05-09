package value

import "fmt"

type Boolean struct {
	Value bool
	proto *Proto
}

func (v *Boolean) v() vtype {
	return value
}

func (v *Boolean) Hash() string {
	return fmt.Sprintf("b:%t", v.Value)
}

func (v *Boolean) Proto() *Proto {
	return v.proto
}

func (v *Boolean) Inherit(p *Proto) Value {
	b := NewBoolean(v.Value)
	b.proto = p

	return b
}

func NewBoolean(v bool) *Boolean {
	return &Boolean{
		Value: v,
		proto: ProtoBoolean,
	}
}

var ProtoBoolean = &Proto{
	Methods: map[string]Caller{},
}

// Compile time checks
var _ Value = (*Boolean)(nil)
