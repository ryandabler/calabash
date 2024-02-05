package value_test

import (
	"calabash/ast"
	"calabash/internal/tokentype"
	"calabash/internal/value"
	"calabash/lexer/tokens"
	"testing"
)

func TestFunctionArity(t *testing.T) {
	t.Run("empty param list", func(t *testing.T) {
		f := &value.Function{}

		if f.Arity() != 0 {
			t.Error("empty function paramters should have arity of zero")
		}
	})

	t.Run("non-rest, non-empty param list", func(t *testing.T) {
		f := &value.Function{
			ParamList: []ast.Identifier{
				{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
				{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 0)},
			},
		}

		if f.Arity() != 2 {
			t.Error("2-parameter list should count as 2-arity")
		}
	})

	t.Run("rest but otherwise empty param list", func(t *testing.T) {
		f := &value.Function{
			ParamList: []ast.Identifier{
				{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0), Rest: true},
			},
		}

		if f.Arity() != 0 {
			t.Error("rest param should not count towards arity")
		}
	})

	t.Run("rest and non empty param list", func(t *testing.T) {
		f := &value.Function{
			ParamList: []ast.Identifier{
				{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
				{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 0), Rest: true},
			},
		}

		if f.Arity() != 1 {
			t.Error("only non-rest param should count towards arity")
		}
	})
}
