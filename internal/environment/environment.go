package environment

type Environment[T any] struct {
	Fields map[string]T
	Parent *Environment[T]
}

func (e *Environment[T]) Add(k string, v T) {
	e.Fields[k] = v
}

func (e *Environment[T]) Get(k string) T {
	v, ok := e.Fields[k]

	if !ok && e.Parent != nil {
		return e.Parent.Get(k)
	}

	return v
}

func (e *Environment[T]) Set(k string, v T) {
	_, ok := e.Fields[k]

	if ok {
		e.Fields[k] = v
		return
	}

	if e.Parent != nil {
		e.Parent.Set(k, v)
	}
}

func (e *Environment[T]) Has(k string) bool {
	_, ok := e.Fields[k]

	if !ok && e.Parent != nil {
		return e.Parent.Has(k)
	}

	return ok
}

func (e *Environment[T]) HasDirectly(k string) bool {
	_, ok := e.Fields[k]

	return ok
}

func New[T any](e *Environment[T]) *Environment[T] {
	return &Environment[T]{Fields: make(map[string]T), Parent: e}
}

func Slice[T any](e *Environment[T], l uint64) *Environment[T] {
	if l == 0 {
		return nil
	}

	env := New[T](nil)
	env.Fields = e.Fields
	curNew := env
	curOld := e

	// Iterate down environment chain, linking `e`'s parents to
	// `env`'s until we either hit depth or exhaust sliced
	// environment's depth
	for n := uint64(1); n < l && curOld.Parent != nil; n++ {
		curNew.Parent = curOld.Parent
		curNew = curNew.Parent
		curOld = curOld.Parent
	}

	// Break chain in new environment so that slice blocks off
	// remaining environment layers
	curNew.Parent = nil

	return env
}
