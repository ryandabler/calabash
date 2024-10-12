package value

import (
	"calabash/ast"
	"calabash/internal/environment"
	"calabash/internal/slice"
	"calabash/internal/uuid"
	"calabash/lexer/tokens"
	"strconv"
)

type ProtoMethod struct {
	ParamList []ast.Identifier
	Apps      []Value
	Depth     struct {
		Specified bool
		Tk        *tokens.Token
	}
	Me          Value
	hash        string
	call        func(me Value, e Evaluator) (interface{}, error)
	Inheritable *Proto
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
		Me:          me,
		ParamList:   pm.ParamList,
		Depth:       pm.Depth,
		Apps:        pm.Apps,
		call:        pm.call,
		hash:        pm.hash,
		Inheritable: pm.Inheritable,
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

// TODO: DRY up with Function#Closure()
func (pm *ProtoMethod) Closure(env *environment.Environment[Value]) *environment.Environment[Value] {
	// Case for `fn () ...` declarations: no closure
	if !pm.Depth.Specified {
		return nil
	}

	// Case for `fn<> () ...` declaractions; full exposure
	if pm.Depth.Tk == nil {
		return env
	}

	// Case for `fn<#> () ...` declarations; limited exposure
	lex := pm.Depth.Tk.Lexeme
	d, _ := strconv.ParseUint(lex, 10, 64) // ignore error since static analyzer will catch these

	return environment.Slice(env, d)
}

func ProtoMethodFromFn(fn *Function) *ProtoMethod {
	return &ProtoMethod{
		Me:        nil,
		ParamList: fn.ParamList,
		Depth:     fn.Depth,
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
