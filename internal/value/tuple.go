package value

import (
	"calabash/internal/slice"
	"fmt"
)

type Tuple struct {
	Items []Value
	proto *Proto
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

// Compile time checks
var _ Value = (*Tuple)(nil)
