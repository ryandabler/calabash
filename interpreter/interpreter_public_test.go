package interpreter_test

import (
	"calabash/ast"
	"calabash/internal/tokentype"
	"calabash/internal/value"
	"calabash/interpreter"
	"calabash/lexer/scanner"
	"calabash/lexer/tokens"
	"calabash/parser"
	staticanalyzer "calabash/static_analyzer"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	t.Run("correct execution", func(t *testing.T) {
		table := []struct {
			name     string
			text     string
			validate func(interface{}, interpreter.IntpState) error
		}{
			{
				name: "literal string 1",
				text: "'abcd'",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewString("abcd")) {
						return errors.New("Values does not equal \"abcd\"")
					}

					return nil
				},
			},
			{
				name: "literal string 2",
				text: "\"abcd\"",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewString("abcd")) {
						return errors.New("Values does not equal \"abcd\"")
					}

					return nil
				},
			},
			{
				name: "literal number 1",
				text: "123",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(123)) {
						return errors.New("Values does not equal 123")
					}

					return nil
				},
			},
			{
				name: "literal number 2",
				text: "123.4",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(123.4)) {
						return errors.New("Values does not equal 123.4")
					}

					return nil
				},
			},
			{
				name: "literal boolean 1",
				text: "true",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Values does not equal true")
					}

					return nil
				},
			},
			{
				name: "literal boolean 2",
				text: "false",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Values does not equal false")
					}

					return nil
				},
			},
			{
				name: "literal bottom",
				text: "bottom",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, &value.Bottom{}) {
						return errors.New("Values does not equal bottom")
					}

					return nil
				},
			},
			{
				name: "literal function 1",
				text: "fn (a) {}",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					vfunc, ok := v.(*value.Function)

					if !ok {
						return errors.New("Did not receive a function value")
					}

					if !reflect.DeepEqual(vfunc.ParamList, []ast.Identifier{{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 4)}}) {
						return errors.New("Param list is not equal")
					}

					if !reflect.DeepEqual(vfunc.Body, ast.Block{Contents: []ast.Node{}}) {
						return errors.New("Function bodies are not equal")
					}

					if !reflect.DeepEqual(vfunc.Apps, []value.Value(nil)) {
						return errors.New("Partial applications are not equal")
					}

					return nil
				},
			},
			{
				name: "literal function 2",
				text: "fn (a, b) -> 1 + 2",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					vfunc, ok := v.(*value.Function)

					if !ok {
						return errors.New("Did not receive a function value")
					}

					if !reflect.DeepEqual(vfunc.ParamList, []ast.Identifier{
						{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 4)},
						{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 7)},
					}) {
						return errors.New("Param list is not equal")
					}

					if !reflect.DeepEqual(vfunc.Body, ast.Block{
						Contents: []ast.Node{
							ast.ReturnStmt{
								Expr: ast.BinaryExpr{
									Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 13)},
									Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 17)},
									Operator: tokens.New(tokentype.PLUS, "+", 0, 15),
								},
							},
						},
					}) {
						return errors.New("Function bodies are not equal")
					}

					if !reflect.DeepEqual(vfunc.Apps, []value.Value(nil)) {
						return errors.New("Partial applications are not equal")
					}

					return nil
				},
			},
			{
				name: "literal tuple 1",
				text: "[1, \"a\", fn() -> 1]",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					tuple, ok := v.(*value.Tuple)

					if !ok {
						return errors.New("Literal value was not a tuple")
					}

					if !reflect.DeepEqual(tuple.Items[0], value.NewNumber(1)) {
						return errors.New("First tuple item is not equal to 1")
					}

					if !reflect.DeepEqual(tuple.Items[1], value.NewString("a")) {
						return errors.New("Second tuple item is not equal to \"a\"")
					}

					fn, ok := tuple.Items[2].(*value.Function)

					if !ok {
						return errors.New("Third tuple item is not a function")
					}

					if !reflect.DeepEqual(
						fn.Body,
						ast.Block{Contents: []ast.Node{ast.ReturnStmt{Expr: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 17)}}}},
					) {
						return errors.New("Function bodies are not the same")
					}

					return nil
				},
			},
			{
				name: "literal tuple 2",
				text: "[]",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					tpl, ok := v.(*value.Tuple)

					if !ok {
						return errors.New("Did not receive a tuple")
					}

					if len(tpl.Items) != 0 {
						return errors.New("Tuple should be empty")
					}

					return nil
				},
			},
			{
				name: "literal record 1",
				text: "{}",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					rec, ok := v.(*value.Record)

					if !ok {
						return errors.New("Did not receive a record")
					}

					if len(rec.Entries) > 0 {
						return errors.New("Record should be empty")
					}

					return nil
				},
			},
			{
				name: "literal record 2",
				text: "{ 'a' -> 1, bottom -> 2, 5 -> 3, true -> 4, [1] -> 5, {} -> 6, (1 + 9) -> 7 }",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					rec, ok := v.(*value.Record)

					if !ok {
						return errors.New("Did not receive a record")
					}

					if len(rec.Entries) != 7 {
						return errors.New("Record should have seven properties")
					}

					v1 := value.NewString("a")
					v2 := &value.Bottom{}
					v3 := value.NewNumber(5)
					v4 := value.NewBoolean(true)
					v5 := value.NewTuple([]value.Value{value.NewNumber(1)})
					v6 := value.NewRecord([]struct {
						K value.Value
						V value.Value
					}{})
					v7 := value.NewNumber(10)

					if v, ok := rec.Entries[v1.Hash()]; !ok {
						return errors.New("Record does not contain key 'a'")
					} else if !reflect.DeepEqual(v, value.NewNumber(1)) {
						return errors.New("Record property 'a' was not assigned the value 1")
					}

					if v, ok := rec.Entries[v2.Hash()]; !ok {
						return errors.New("Record does not contain key bottom")
					} else if !reflect.DeepEqual(v, value.NewNumber(2)) {
						return errors.New("Record property bottom was not assigned the value 2")
					}

					if v, ok := rec.Entries[v3.Hash()]; !ok {
						return errors.New("Record does not contain key 5")
					} else if !reflect.DeepEqual(v, value.NewNumber(3)) {
						return errors.New("Record property 5 was not assigned the value 3")
					}

					if v, ok := rec.Entries[v4.Hash()]; !ok {
						return errors.New("Record does not contain key true")
					} else if !reflect.DeepEqual(v, value.NewNumber(4)) {
						return errors.New("Record property true was not assigned the value 4")
					}

					if v, ok := rec.Entries[v5.Hash()]; !ok {
						return errors.New("Record does not contain key [1]")
					} else if !reflect.DeepEqual(v, value.NewNumber(5)) {
						return errors.New("Record property [1] was not assigned the value 5")
					}
					if v, ok := rec.Entries[v6.Hash()]; !ok {
						return errors.New("Record does not contain key {}")
					} else if !reflect.DeepEqual(v, value.NewNumber(6)) {
						return errors.New("Record property {} was not assigned the value 6")
					}
					if v, ok := rec.Entries[v7.Hash()]; !ok {
						return errors.New("Record does not contain key 10")
					} else if !reflect.DeepEqual(v, value.NewNumber(7)) {
						return errors.New("Record property 10 was not assigned the value 7")
					}

					return nil
				},
			},
			{
				name: "literal record 3",
				text: "let a = fn () -> {}; { a -> 1}",
				validate: func(v interface{}, is interpreter.IntpState) error {
					rec, ok := v.(*value.Record)

					if !ok {
						return errors.New("Did not receive a record")
					}

					if len(rec.Entries) != 1 {
						return errors.New("Record should have seven properties")
					}

					fn := is.Env.Get("a")

					val, ok := rec.Entries[fn.Hash()]

					if !ok {
						return errors.New("Function was not properly keyed to object")
					}

					if !reflect.DeepEqual(val, value.NewNumber(1)) {
						return errors.New("Value for function key is not the number 1")
					}

					return nil
				},
			},
			{
				name: "literal record proto 'get'",
				text: "{ 'a' -> 1 }->'get'('a')",
				validate: func(v interface{}, is interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(1)) {
						return errors.New("Did not properly get value for key 'a'")
					}

					return nil
				},
			},
			{
				name: "binary addition 1",
				text: "1 + 1",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(2)) {
						return errors.New("Values does not equal 2")
					}

					return nil
				},
			},
			{
				name: "binary addition 2",
				text: "'1' + '1'",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewString("11")) {
						return errors.New("Values does not equal '11'")
					}

					return nil
				},
			},
			{
				name: "binary subtraction",
				text: "2 - 5",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(-3)) {
						return errors.New("Values does not equal -3")
					}

					return nil
				},
			},
			{
				name: "binary multiplication",
				text: "2 * 5",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(10)) {
						return errors.New("Values does not equal 10")
					}

					return nil
				},
			},
			{
				name: "binary division",
				text: "5 / 3",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(5.0/3.0)) {
						return errors.New(fmt.Sprintf("Values does not equal %f", 5.0/3.0))
					}

					return nil
				},
			},
			{
				name: "binary exponentiation",
				text: "5 ** 2",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(25)) {
						return errors.New("Values does not equal 25")
					}

					return nil
				},
			},
			{
				name: "binary greater",
				text: "5 > 2",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Values does not equal true")
					}

					return nil
				},
			},
			{
				name: "binary greater or equal",
				text: "5 >= 2",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Values does not equal true")
					}

					return nil
				},
			},
			{
				name: "binary lesser",
				text: "5 < 2",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Values does not equal false")
					}

					return nil
				},
			},
			{
				name: "binary lesser or equal",
				text: "5 <= 2",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Values does not equal false")
					}

					return nil
				},
			},
			{
				name: "binary equal to",
				text: "2 == 3",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Values does not equal false")
					}

					return nil
				},
			},
			{
				name: "binary equal to (tuples) 1",
				text: "[1, \"a\", [true, bottom]] == [1, \"a\", [true, bottom]]",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Tuples should be deeply equal")
					}

					return nil
				},
			},
			{
				name: "binary equal to (tuples) 2",
				text: "[1, \"a\", [true, bottom]] == [\"a\", [true, bottom], 1]",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Out of order tuples should not be deeply equal")
					}

					return nil
				},
			},
			{
				name: "binary equal to (tuples) 3",
				text: "[1, \"a\", [true, bottom]] == [\"a\", [true, bottom]]",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Tuples with different numbers of elements should not be deeply equal")
					}

					return nil
				},
			},
			{
				name: "binary not equal to",
				text: "2 != 3",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Values does not equal true")
					}

					return nil
				},
			},
			{
				name: "binary boolean and",
				text: "true && false",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Values does not equal false")
					}

					return nil
				},
			},
			{
				name: "binary boolean and short circuit",
				text: "false && 1",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Values does not equal false")
					}

					return nil
				},
			},
			{
				name: "binary boolean or",
				text: "true || false",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Values does not equal true")
					}

					return nil
				},
			},
			{
				name: "binary boolean or short circuit",
				text: "true || 1",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(true)) {
						return errors.New("Values does not equal true")
					}

					return nil
				},
			},
			{
				name: "grouping expression",
				text: "(5 ** 2)",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(25)) {
						return errors.New("Values does not equal 25")
					}

					return nil
				},
			},
			{
				name: "pipe expression 1",
				text: "1 |> ? + 1",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(2)) {
						return errors.New("Pipe expression did not properly resolve value for '?'")
					}

					return nil
				},
			},
			{
				name: "pipe expression 2",
				text: "1 |> ? + 1 |> 3 ** ?",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(9)) {
						return errors.New("Pipe expression did not properly resolve value for chains of '?'")
					}

					return nil
				},
			},
			{
				name: "pipe expression can properly nest",
				text: "1 |> ? + (2 |> ? - 1) |> 3 * ?",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(6)) {
						return errors.New("Could not properly resolve nested '?'")
					}

					return nil
				},
			},
			{
				name: "container values can be composed from '?'",
				text: "1 |> [?, ? + 1]",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewTuple([]value.Value{value.NewNumber(1), value.NewNumber(2)})) {
						return errors.New("Could not properly build tuple from '?'")
					}

					return nil
				},
			},
			{
				name: "question mark does not remain after pipe expression",
				text: "1 |> ? ",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					q := i.Env.Get("?")

					if q != nil {
						return errors.New("? outlived the pipe expression")
					}

					return nil
				},
			},
			{
				name: "unary minus",
				text: "-5",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(-5)) {
						return errors.New("Values does not equal -5")
					}

					return nil
				},
			},
			{
				name: "variable declaration (no init)",
				text: "let a;",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !i.Env.HasDirectly("a") {
						return errors.New("Variable 'a' should be present in first layer of state")
					}

					if !reflect.DeepEqual(i.Env.Get("a"), &value.Bottom{}) {
						return errors.New("Undefined variables should be initialized with bottom value")
					}

					return nil
				},
			},
			{
				name: "multi variable declaration (no init)",
				text: "let a, b;",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !i.Env.HasDirectly("a") {
						return errors.New("Variable 'a' should be present in first layer of state")
					}

					if !i.Env.HasDirectly("b") {
						return errors.New("Variable 'b' should be present in first layer of state")
					}

					if !reflect.DeepEqual(i.Env.Get("a"), &value.Bottom{}) && !reflect.DeepEqual(i.Env.Get("b"), &value.Bottom{}) {
						return errors.New("Undefined variables should be initialized with bottom value")
					}

					return nil
				},
			},
			{
				name: "multi variable declaration (with init) 1",
				text: "let a, b = 1, 3 + 4;",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("Variable 'a' was not set to value 1")
					}

					if !reflect.DeepEqual(i.Env.Get("b"), value.NewNumber(7)) {
						return errors.New("Variable 'b' was not resolved to value 7")
					}

					return nil
				},
			},
			{
				name: "multi variable declaration (with init) 2",
				text: "let a, b = 1, a;",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					n := value.NewNumber(1)

					if !reflect.DeepEqual(i.Env.Get("a"), n) {
						return errors.New("Variable 'a' was not set to value 1")
					}

					if !reflect.DeepEqual(i.Env.Get("b"), n) {
						return errors.New("Variable 'b' was not resolved to variable \"a\"'s value 1")
					}

					return nil
				},
			},
			{
				name: "identifier expression",
				text: "let a; a",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !i.Env.HasDirectly("a") {
						return errors.New("Variable 'a' should be present in first layer of state")
					}

					if !reflect.DeepEqual(v, &value.Bottom{}) {
						return errors.New("Referencing an undefined variable should give a bottom value")
					}

					return nil
				},
			},
			{
				name: "functions should be called when their arguments list equal their arity",
				text: "let a = fn (a, b) -> a + b; a(1, 2)",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(3)) {
						return errors.New("Function was not properly called")
					}

					return nil
				},
			},
			{
				name: "functions should be partially applied when their arguments list is less than their arity",
				text: "let a = fn (a, b) -> a + b; a(1)",
				validate: func(v interface{}, i interpreter.IntpState) error {
					fn, ok := v.(*value.Function)

					if !ok {
						return errors.New("Partially applied function did not return a function")
					}

					if fn.Arity() != 1 {
						return errors.New("Arity was not updated for a partially applied function")
					}

					if !reflect.DeepEqual(fn.Apps, []value.Value{value.NewNumber(1)}) {
						return errors.New("Function applied arguments were not evaluated properly")
					}

					return nil
				},
			},
			{
				name: "partially applied functions are callable",
				text: "let a = fn (a, b) -> a + b; let b = a(1); b(2)",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(3)) {
						return errors.New("Partially applied arguments were not stored properly")
					}

					return nil
				},
			},
			{
				name: "partially applied functions are not equal",
				text: "let a = fn (a, b) -> a + b; let b = a(1); let c = a(1); b == c",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewBoolean(false)) {
						return errors.New("Partially applied functions with identical applications are not equal")
					}

					return nil
				},
			},
			{
				name: "functions can be applied with more arguments than arity",
				text: "let a = fn (a, b) -> a + b; a(1,2,3)",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(3)) {
						return errors.New("Extra arguments should be discarded when functions are called")
					}

					return nil
				},
			},
			{
				name: "function bodies create their own scope",
				text: "let a, b; let c = fn(mut a) { a = 2; let c = 2; };",
				validate: func(v interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), &value.Bottom{}) {
						return errors.New("Function parameters should be declared in an inner scope")
					}

					if !reflect.DeepEqual(i.Env.Get("b"), &value.Bottom{}) {
						return errors.New("Variables in function bodies should be declared in an inner scope")
					}

					return nil
				},
			},
			{
				name: "iife 1",
				text: "fn(a){ return a + 1; }(1)",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(2)) {
						return errors.New("Function expression was not immediately invoked")
					}

					return nil
				},
			},
			{
				name: "iife 2",
				text: "(fn(a) -> a + 1)(1)",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewNumber(2)) {
						return errors.New("Function expression was not immediately invoked")
					}

					return nil
				},
			},
			{
				name: "proto expression",
				text: "proto { 'a' -> fn() -> 1 + 1 }",
				validate: func(v interface{}, is interpreter.IntpState) error {
					p, ok := v.(*value.Proto)

					if !ok {
						return errors.New("Did not receive a proto value")
					}

					if len(p.Methods) != 1 {
						return errors.New("Proto method set is not of size 1")
					}

					_, ok = p.Methods["s:\"a\""]

					if !ok {
						return errors.New("Method in method set is not keyed by 'a'")
					}

					return nil
				},
			},
			{
				name: "protos can be bound to different instances of same type",
				text: "let a, b = []->'push', []->'push'; [a, b]",
				validate: func(v interface{}, is interpreter.IntpState) error {
					a := is.Env.Get("a")
					b := is.Env.Get("b")

					apm, ok := a.(*value.ProtoMethod)

					if !ok {
						return errors.New("Did not receive a proto method for variable 'a'")
					}

					bpm, ok := b.(*value.ProtoMethod)

					if !ok {
						return errors.New("Did not receive a proto method for variable 'b'")
					}

					if apm.Me == bpm.Me {
						return errors.New("Each proto method must point to the different instances")
					}

					return nil
				},
			},
			{
				name: "protos can be bound to different instances of same type",
				text: "let t, a, b = [], t->'push', t->'push';",
				validate: func(_ interface{}, is interpreter.IntpState) error {
					a := is.Env.Get("a")
					b := is.Env.Get("b")

					apm, ok := a.(*value.ProtoMethod)

					if !ok {
						return errors.New("Did not receive a proto method for variable 'a'")
					}

					bpm, ok := b.(*value.ProtoMethod)

					if !ok {
						return errors.New("Did not receive a proto method for variable 'b'")
					}

					if apm.Me != bpm.Me {
						return errors.New("Each proto method must point to the same instance")
					}

					return nil
				},
			},
			{
				name: "proto method tests: tuple -> 'push'",
				text: "let a, b = [], a->'push'(1);",
				validate: func(_ interface{}, is interpreter.IntpState) error {
					a := is.Env.Get("a")
					b := is.Env.Get("b")

					atpl, ok := a.(*value.Tuple)

					if !ok {
						return errors.New("Did not receive a tuple for variable 'a'")
					}

					btpl, ok := b.(*value.Tuple)

					if !ok {
						return errors.New("Did not receive a tuple for variable 'b'")
					}

					if len(atpl.Items) != 0 {
						return errors.New("Tuple was mutated inadvertently")
					}

					if len(btpl.Items) != 1 {
						return errors.New("Tuple was incorrectly appended to")
					}

					return nil
				},
			},
			{
				name: "proto method tests: tuple -> 'push'",
				text: "let a, b, c = [], a->'push'(1), a->'push'(2);",
				validate: func(_ interface{}, is interpreter.IntpState) error {
					b := is.Env.Get("b")
					c := is.Env.Get("c")

					btpl, ok := b.(*value.Tuple)

					if !ok {
						return errors.New("Did not receive a tuple for variable 'b'")
					}

					ctpl, ok := c.(*value.Tuple)

					if !ok {
						return errors.New("Did not receive a tuple for variable 'c'")
					}

					if len(btpl.Items) != 1 || len(ctpl.Items) != 1 {
						return errors.New("Tuples were not given the right number of arguments")
					}

					n := value.NewNumber(1)

					if !reflect.DeepEqual(btpl.Items[0], n) {
						return errors.New("Tuple 'b' should only have value 1 inside")
					}

					n.Value = 2

					if !reflect.DeepEqual(ctpl.Items[0], n) {
						return errors.New("Tuple 'b' should only have value 2 inside")
					}

					return nil
				},
			},
			{
				name: "proto gets can be chained with calls",
				text: "[] -> 'push'(1) -> 'push'(2)",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					tpl, ok := v.(*value.Tuple)

					if !ok {
						return errors.New("Did not receive a tuple")
					}

					if len(tpl.Items) != 2 {
						return errors.New("Tuple should have exactly two elements")
					}

					n := value.NewNumber(1)

					if !reflect.DeepEqual(tpl.Items[0], n) {
						return errors.New("First element should be 1 for tuple")
					}

					n.Value = 2

					if !reflect.DeepEqual(tpl.Items[1], n) {
						return errors.New("Second element should be 2 for tuple")
					}

					return nil
				},
			},
			{
				name: "proto method tests: Number->'stringify'",
				text: "1->'stringify'()",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewString("1")) {
						return errors.New("Stringified numeric value should be '1'")
					}

					return nil
				},
			},
			{
				name: "proto method tests: String->'upper'",
				text: "'abc'->'upper'()",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewString("ABC")) {
						return errors.New("Uppercased value should be 'ABC'")
					}

					return nil
				},
			},
			{
				name: "proto method tests: Boolean->'stringify'",
				text: "true->'stringify'()",
				validate: func(v interface{}, _ interpreter.IntpState) error {
					if !reflect.DeepEqual(v, value.NewString("true")) {
						return errors.New("Stringified boolean value should be 'true'")
					}

					return nil
				},
			},
			{
				name: "assign statement 1",
				text: "let mut a; a = 4;",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(4)) {
						return errors.New("Variable \"a\" was not properly assigned the value 4")
					}

					return nil
				},
			},
			// {
			// 	name: "assign statement 2",
			// 	text: "let mut a, mut b = 1, 2; a, b = b, a;",
			// 	validate: func(_ interface{}, i interpreter.IntpState) error {
			// 		if i.Env.Get("a") != (value.VNumber{Value: 2}) {
			// 			return errors.New("Variable \"a\" was not properly assigned \"b\"'s value 2")
			// 		}

			// 		if i.Env.Get("b") != (value.VNumber{Value: 1}) {
			// 			return errors.New("Variable \"b\" was not properly assigned \"a\"'s value 1")
			// 		}

			// 		return nil
			// 	},
			// },
			{
				name: "if statement (no init) enters then block",
				text: "let mut a; if true { a = 1; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("The `then` block in the if statement was not entered for a true value")
					}

					return nil
				},
			},
			{
				name: "if statement (no init) enters else block",
				text: "let mut a; if false { a = 1; } else { a = 2; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(2)) {
						return errors.New("The `else` block in the if statement was not entered for a false value")
					}

					return nil
				},
			},
			{
				name: "if statement (no init) does nothing if false and no else block",
				text: "let mut a; if false { a = 1; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), &value.Bottom{}) {
						return errors.New("Variable \"a\" should not have been reassigned")
					}

					return nil
				},
			},
			{
				name: "if statement (no init) generates new scope for each block",
				text: "let a = 1; if true { let a = 2; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("Outer variable \"a\" should not have been reassigned")
					}

					return nil
				},
			},
			{
				name: "if statement (with init) shadows outer variables",
				text: "let a = 1; if let mut a = 2; true { a = 3; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("Outer variable \"a\" should have been shadowed")
					}

					return nil
				},
			},
			{
				name: "else statement can reference initialized if variables",
				text: "let mut b; if let a = 2; false {} else { b = a; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("b"), value.NewNumber(2)) {
						return errors.New("Else branches should be able to access variables declared in `if` blocks")
					}

					return nil
				},
			},
			{
				name: "nested if statements should be able to access outer if statements variables",
				text: "let mut b; if let a = 2; false {} else if true { b = a; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("b"), value.NewNumber(2)) {
						return errors.New("Nested if statements should be able to access outer if statements' variables")
					}

					return nil
				},
			},
			{
				name: "nested if statements should be able to shadow outer if statements variables",
				text: "let mut b; if let a = 2; false {} else if let a = 3; true { b = a; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("b"), value.NewNumber(3)) {
						return errors.New("Nested if statements should be able to shadow outer if statements' variables")
					}

					return nil
				},
			},
			{
				name: "while statements create their own scope",
				text: "let a = 1; while let mut a = 10; a < 11 { a = a + 1; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("Variables initialized in while loop should shadow outer scope")
					}

					return nil
				},
			},
			{
				name: "while loops do not execute block when condition is false",
				text: "let mut a = 1; while false { a = a + 1; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("While loop body should not have been entered")
					}

					return nil
				},
			},
			{
				name: "while loops can be broken out of",
				text: "let mut a = 1; while true { if a == 1 { break; } a = a + 1; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("a"), value.NewNumber(1)) {
						return errors.New("While loop body was not properly broken out of")
					}

					return nil
				},
			},
			{
				name: "while loops can be continued",
				text: "let mut a, mut b = 1, 10; while a <= 3 { a = a + 1;  if a == 2 { continue; } b = b + 1; }",
				validate: func(_ interface{}, i interpreter.IntpState) error {
					if !reflect.DeepEqual(i.Env.Get("b"), value.NewNumber(12)) {
						return errors.New("While loop body was not properly continued")
					}

					return nil
				},
			},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: Unexpected error during lexing phase", e.name)
				continue
			}

			ast, err := parser.New(ts).Parse()

			if err != nil {
				t.Errorf("%q: Unexpected error parsing text", e.name)
				continue
			}

			err = staticanalyzer.New().Analyze(ast)

			if err != nil {
				t.Errorf("%q: Unexpected static error %q", e.name, err)
				continue
			}

			i := interpreter.New()
			v, err := i.Eval(ast)

			if err != nil {
				t.Errorf("%q: Unexpected runtime error %q", e.name, err)
			}

			err = e.validate(v, i.Dump())

			if err != nil {
				t.Errorf("%q: %s", e.name, err)
			}
		}
	})

	t.Run("runtime errors", func(t *testing.T) {
		table := []struct {
			name string
			text string
		}{
			{
				name: "Only boolean values can be used in '&&' expressions",
				text: "1 && true",
			},
			{
				name: "Only boolean values can be used in '||' expressions",
				text: "false || 'a'",
			},
			{
				name: "Types cannot be different for binary addition",
				text: "'a' + 1",
			},
			{
				name: "Numeric binary operators require only numbers: mixed with string",
				text: "'a' - 1",
			},
			{
				name: "Numeric binary operators require only numbers: mixed with boolean",
				text: "true - 1",
			},
			{
				name: "Numeric binary operators require only numbers: mixed with bottom",
				text: "bottom - 1",
			},
			{
				name: "Numeric binary operators require only numbers: mixed with function",
				text: "(fn() -> 1) - 1",
			},
			{
				name: "Unary minus operator only works with numbers",
				text: "-false",
			},
			{
				name: "Non-functions are not callable",
				text: "let a = 1; a()",
			},
			{
				name: "Errors in function calls bubble up",
				text: "fn() { 1 + 'a' }()",
			},
			{
				name: "Errors in else blocks bubble up",
				text: "if true { 1 + 'b' }",
			},
			{
				name: "Errors in if declarations should bubble up",
				text: "if let a = 1 + '1'; true {}",
			},
			{
				name: "If conditions must be boolean",
				text: "if 1 {}",
			},
			{
				name: "If conditions should bubble errors",
				text: "if 1 + '1' {}",
			},
			{
				name: "Else blocks should bubble errors",
				text: "if false {} else { 1 + '1' }",
			},
			{
				name: "Errors in return statements should bubble up",
				text: "fn () { return 1 + '1'; }()",
			},
			{
				name: "Errors in lhs of a binary expression bubble up",
				text: "(1 + '1') + 2",
			},
			{
				name: "Errors in rhs of a binary expression bubble up",
				text: "2 + (1 + '1')",
			},
			{
				name: "Errors in unary expressions bubble up",
				text: "-(1 + '1')",
			},
			{
				name: "Errors in resolving callee bubble up",
				text: "(1 + '1')()",
			},
			{
				name: "Errors in a function's arguments list bubble up",
				text: "fn (a, b) {}(1, 1 + '1')",
			},
			{
				name: "Errors in assignments bubble up",
				text: "let mut a, mut b; a, b = 1, 1 + '1';",
			},
			{
				name: "Errors in tuples should bubble up",
				text: "[1 + 'a']",
			},
			{
				name: "bottom value has no proto",
				text: "bottom -> 'a'",
			},
			{
				name: "function value has no proto",
				text: "let a = fn () -> bottom; a -> 'a'",
			},
			{
				name: "proto value has no proto",
				text: "let p = proto { 'a' -> fn () -> bottom }; p -> 'a'",
			},
			{
				name: "protoMethod value has no proto",
				text: "let a = [] -> 'push'; a -> 'a'",
			},
			{
				name: "while loop with non-boolean condition expression 1",
				text: "while 1 { 1 + 1 }",
			},
			{
				name: "while loop with non-boolean condition expression 2",
				text: "while 'true' { 1 + 1 }",
			},
			{
				name: "getting record property that doesn't exist",
				text: "{ 'a' -> 1 }->'get'(1)",
			},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: Unexpected error during lexing phase", e.name)
				continue
			}

			ast, err := parser.New(ts).Parse()

			if err != nil {
				t.Errorf("%q: Unexpected error parsing text", e.name)
				continue
			}

			err = staticanalyzer.New().Analyze(ast)

			if err != nil {
				t.Errorf("%q: Unexpected static error %q", e.name, err)
				continue
			}

			i := interpreter.New()
			_, err = i.Eval(ast)

			if err == nil {
				t.Errorf("%q: Unexpected runtime success", e.name)
			}
		}
	})
}
