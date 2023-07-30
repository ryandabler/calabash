package value

import (
	"calabash/ast"
	"calabash/errors"
	"calabash/internal/tokentype"
	"calabash/lexer/tokens"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	ProtoTuple.Methods["s:\"push\""] = &ProtoMethod{
		ParamList: []ast.Identifier{
			{Name: tokens.New(tokentype.IDENTIFIER, "e", 0, 0), Mut: false},
		},
		call: func(me Value, e Evaluator) (interface{}, error) {
			tpl, ok := me.(*Tuple)

			if !ok {
				return nil, errors.RuntimeError{Msg: "Cannot call `push` on a value other than a tuple"}
			}

			vs := append(tpl.Items, e.Dump().Env.Get("e"))

			return NewTuple(vs), nil
		},
	}

	ProtoNumber.Methods["s:\"stringify\""] = &ProtoMethod{
		call: func(me Value, _ Evaluator) (interface{}, error) {
			n, ok := me.(*Number)

			if !ok {
				return nil, errors.RuntimeError{Msg: "Expect 'me' to be a number"}
			}

			s := strconv.FormatFloat(float64(n.Value), 'f', -1, 64)

			return NewString(s), nil
		},
	}

	ProtoString.Methods["s:\"upper\""] = &ProtoMethod{
		call: func(me Value, _ Evaluator) (interface{}, error) {
			s, ok := me.(*String)

			if !ok {
				return nil, errors.RuntimeError{Msg: "Expect 'me' to be a string"}
			}

			return NewString(strings.ToUpper(s.Value)), nil
		},
	}

	ProtoBoolean.Methods["s:\"stringify\""] = &ProtoMethod{
		call: func(me Value, _ Evaluator) (interface{}, error) {
			n, ok := me.(*Boolean)

			if !ok {
				return nil, errors.RuntimeError{Msg: "Expect 'me' to be a boolean"}
			}

			return NewString(fmt.Sprint(n.Value)), nil
		},
	}
}
