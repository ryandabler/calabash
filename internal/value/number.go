package value

import "fmt"

type Number struct {
	Value float64
	proto *Proto
}

func (v *Number) v() vtype {
	return value
}

func (v *Number) Hash() string {
	return fmt.Sprintf("n:%v", v.Value)
}

func (v *Number) Proto() *Proto {
	return v.proto
}

func NewNumber(v float64) *Number {
	return &Number{
		Value: v,
		proto: ProtoNumber,
	}
}

var ProtoNumber = &Proto{
	Methods: map[string]Caller{},
}

// Compile time checks
var _ Value = (*Number)(nil)
