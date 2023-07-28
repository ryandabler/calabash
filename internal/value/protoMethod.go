package value

import (
	"calabash/ast"
	"calabash/internal/uuid"
)

type ProtoMethod struct {
	ParamList []ast.Identifier
	Apps      []Value
	Me        Value
	hash      string
	call      func(me Value, e Evaluator) (interface{}, error)
}

func (pm *ProtoMethod) v() vtype {
	return value
}

func (pm *ProtoMethod) Proto() *Proto {
	return nil
}

func (pm *ProtoMethod) Apply(vs []Value) Caller {
	return &ProtoMethod{
		ParamList: pm.ParamList,
		Apps:      append(pm.Apps, vs...),
		Me:        pm.Me,
		call:      pm.call,
	}
}

func (pm *ProtoMethod) Bind(me Value) Caller {
	if pm.Me != nil {
		return pm
	}

	return &ProtoMethod{
		Me:        me,
		ParamList: pm.ParamList,
		Apps:      pm.Apps,
		call:      pm.call,
		hash:      pm.hash,
	}
}

func (pm *ProtoMethod) Args() []Value {
	return pm.Apps
}

func (pm *ProtoMethod) Params() []ast.Identifier {
	return pm.ParamList
}

func (pm *ProtoMethod) Arity() int {
	return len(pm.ParamList) - len(pm.Apps)
}

func (pm *ProtoMethod) Call(e Evaluator) (interface{}, error) {
	return pm.call(pm.Me, e)
}

func (pm *ProtoMethod) Hash() string {
	if pm.hash == "" {
		pm.hash = uuid.V4()
	}

	return pm.hash
}

// Compile time check
var _ Value = (*ProtoMethod)(nil)
var _ Caller = (*ProtoMethod)(nil)
