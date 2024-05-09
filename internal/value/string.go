package value

import "fmt"

type String struct {
	Value string
	proto *Proto
}

func (v *String) v() vtype {
	return value
}

func (v *String) Hash() string {
	return fmt.Sprintf("s:%q", v.Value)
}

func (v *String) Proto() *Proto {
	return v.proto
}

func (v *String) Inherit(p *Proto) Value {
	s := NewString(v.Value)
	s.proto = p

	return s
}

func NewString(v string) *String {
	return &String{
		Value: v,
		proto: ProtoString,
	}
}

var ProtoString = &Proto{
	Methods: map[string]Caller{},
}

// Compile time checks
var _ Value = (*String)(nil)
