package environment

import (
	"calabash/internal/value"
)

type Environment struct {
	Fields map[string]value.Value
	Parent *Environment
}

func (e *Environment) Add(k string, v value.Value) {
	e.Fields[k] = v
}

func (e *Environment) Get(k string) value.Value {
	v, ok := e.Fields[k]

	if !ok && e.Parent != nil {
		return e.Parent.Get(k)
	}

	return v
}

func (e *Environment) Set(k string, v value.Value) {
	_, ok := e.Fields[k]

	if !ok {
		e.Parent.Set(k, v)
		return
	}

	e.Fields[k] = v
}

func (e *Environment) Has(k string) bool {
	_, ok := e.Fields[k]

	if !ok && e.Parent != nil {
		return e.Parent.Has(k)
	}

	return ok
}

func (e *Environment) HasDirectly(k string) bool {
	_, ok := e.Fields[k]

	return ok
}

func New(e *Environment) *Environment {
	return &Environment{Fields: make(map[string]value.Value), Parent: e}
}
