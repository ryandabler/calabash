package value

import (
	"calabash/ast"
	"calabash/internal/slice"
	"calabash/internal/uuid"
)

type Function struct {
	ParamList []ast.Identifier
	Body      ast.Block
	Apps      []Value
	hash      string
}

func (v *Function) v() vtype {
	return value
}

func (v *Function) Hash() string {
	if v.hash == "" {
		// Because functions are always unique, to populate the unexported
		// hash we need to manually construct functions in this package
		// with the UUID supplied
		v.hash = uuid.V4()
	}

	return v.hash
}

func (v *Function) Proto() *Proto {
	return nil
}

func (v *Function) Apply(vs []Value) Caller {
	return &Function{
		ParamList: v.ParamList,
		Body:      v.Body,
		Apps:      append(v.Apps, vs...),
	}
}

func (v *Function) Args() []Value {
	return v.Apps
}

func (v *Function) Params() []ast.Identifier {
	return v.ParamList
}

func (v *Function) Rest() bool {
	l, ok := slice.Last(v.ParamList)

	return ok && l.Rest
}

func (v *Function) Arity() int {
	if v.Rest() {
		return len(v.ParamList) - len(v.Apps) - 1
	}

	return len(v.ParamList) - len(v.Apps)
}

func (v *Function) Call(e Evaluator) (interface{}, error) {
	rVal, err := e.Eval(v.Body.Contents)

	if err != nil {
		return nil, err
	}

	if rVal == nil {
		return Bottom{}, nil
	}

	return rVal, nil
}

// Compile time checks
var _ Value = (*Function)(nil)
var _ Caller = (*Function)(nil)
