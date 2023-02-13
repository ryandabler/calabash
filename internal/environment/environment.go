package environment

import (
	"calabash/internal/value"
)

type Environment struct {
	Fields map[string]value.Value
}

func (e *Environment) Add(k string, v value.Value) {
	e.Fields[k] = v
}

func (e *Environment) Get(k string) value.Value {
	v, _ := e.Fields[k]

	return v
}

func (e *Environment) Set(k string, v value.Value) {
	e.Fields[k] = v
}

func (e *Environment) Has(k string) bool {
	_, ok := e.Fields[k]

	return ok
}

func New() *Environment {
	return &Environment{Fields: make(map[string]value.Value)}
}
