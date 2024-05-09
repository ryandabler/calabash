package value

import (
	"calabash/ast"
	"calabash/internal/slice"
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

func (v *ProtoMethod) Inherit(_ *Proto) Value {
	return v
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

func (pm *ProtoMethod) Rest() bool {
	l, ok := slice.Last(pm.ParamList)

	return ok && l.Rest
}

func (pm *ProtoMethod) Arity() int {
	if pm.Rest() {
		return len(pm.ParamList) - len(pm.Apps) - 1
	}

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

func ProtoMethodFromFn(fn *Function) *ProtoMethod {
	return &ProtoMethod{
		Me:        nil,
		ParamList: fn.ParamList,
		Apps:      fn.Apps,
		call: func(me Value, e Evaluator) (interface{}, error) {
			e.PushEnv(nil)
			defer e.PopEnv()

			e.AddEnv("me", me)

			return fn.Call(e)
		},
		hash: fn.hash,
	}
}

// Compile time check
var _ Value = (*ProtoMethod)(nil)
var _ Caller = (*ProtoMethod)(nil)
