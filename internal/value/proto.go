package value

import "calabash/internal/slice"

type Proto struct {
	Methods map[string]Caller
	Keys    []string
	hash    string
}

func (v *Proto) v() vtype {
	return value
}

func (v *Proto) Hash() string {
	if v.hash == "" {
		v.hash = slice.Fold(v.Keys, "prt:", func(k string, acc string, _ int) string {
			return acc + k + "->" + v.Methods[k].Hash()
		})
	}

	return v.hash
}

func (v *Proto) Proto() *Proto {
	return nil
}

// Compile time checks
var _ Value = (*Proto)(nil)
