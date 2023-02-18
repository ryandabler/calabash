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

	if !ok {
		e.Parent.Set(k, v)
		return
	}

	e.Fields[k] = v
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
