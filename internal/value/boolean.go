package value

import "fmt"

type Boolean struct {
	Value bool
}

func (v *Boolean) v() vtype {
	return value
}

func (v *Boolean) Hash() string {
	return fmt.Sprintf("b:%t", v.Value)
}

// Compile time checks
var _ Value = (*Boolean)(nil)
