package value

import "fmt"

type String struct {
	Value string
}

func (v *String) v() vtype {
	return value
}

func (v *String) Hash() string {
	return fmt.Sprintf("s:%q", v.Value)
}

// Compile time checks
var _ Value = (*String)(nil)
