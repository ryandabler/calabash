package value

import (
	"calabash/internal/slice"
	"fmt"
)

type Record struct {
	Entries map[string]Value
	keys    []Value
	proto   *Proto
	hash    string
}

func (v *Record) v() vtype {
	return value
}

func (v *Record) Hash() string {
	if v.hash == "" {
		v.hash = fmt.Sprintf("rec:%s", slice.Fold(v.keys, "", func(k Value, acc string, _ int) string {
			return acc + "," + k.Hash() + ":" + v.Entries[k.Hash()].Hash()
		}))
	}

	return v.hash
}

func (v *Record) Proto() *Proto {
	return v.proto
}

func (v *Record) Inherit(p *Proto) Value {
	es, _ := slice.Map(v.keys, func(k Value) (struct {
		K Value
		V Value
	}, error) {
		return struct {
			K Value
			V Value
		}{K: k, V: v.Entries[k.Hash()]}, nil
	})

	r := NewRecord(es)
	r.proto = p

	return r
}

func NewRecord(vs []struct {
	K Value
	V Value
}) *Record {
	keys, _ := slice.Map(vs, func(v struct {
		K Value
		V Value
	}) (Value, error) {
		return v.K, nil
	})

	return &Record{
		Entries: slice.Fold(keys, map[string]Value{}, func(v Value, acc map[string]Value, i int) map[string]Value {
			acc[v.Hash()] = vs[i].V
			return acc
		}),
		keys:  keys,
		proto: ProtoRecord,
	}
}

var ProtoRecord = &Proto{
	Methods: map[string]Caller{},
}

// Compile time checks
var _ Value = (*Record)(nil)
