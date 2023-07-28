package value

import "fmt"

type Number struct {
	Value float64
}

func (v *Number) v() vtype {
	return value
}

func (v *Number) Hash() string {
	return fmt.Sprintf("n:%v", v.Value)
}

// Compile time checks
var _ Value = (*Number)(nil)
