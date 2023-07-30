package parser_test

import (
	"calabash/ast"
	"calabash/internal/tokentype"
	"calabash/lexer/scanner"
	"calabash/lexer/tokens"
	"calabash/parser"
	"reflect"
	"testing"
)

func nodesAreEqual(a ast.Node, b ast.Node) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	kindA := reflect.TypeOf(a).Kind()
	kindB := reflect.TypeOf(b).Kind()

	if kindA != kindB {
		return false
	}

	tA1, okA := a.(ast.BinaryExpr)
	tB1, okB := b.(ast.BinaryExpr)

	if okA && okB {
		return tA1.Operator.Type == tB1.Operator.Type &&
			nodesAreEqual(tA1.Left, tB1.Left) &&
			nodesAreEqual(tA1.Right, tB1.Right)
	}

	tA2, okA := a.(ast.UnaryExpr)
	tB2, okB := b.(ast.UnaryExpr)

	if okA && okB {
		return tA2.Operator.Type == tB2.Operator.Type &&
			nodesAreEqual(tA2.Expr, tB2.Expr)
	}

	tA3, okA := a.(ast.GroupingExpr)
	tB3, okB := b.(ast.GroupingExpr)

	if okA && okB {
		return nodesAreEqual(tA3.Expr, tB3.Expr)
	}

	tA4, okA := a.(ast.NumericLiteralExpr)
	tB4, okB := b.(ast.NumericLiteralExpr)

	if okA && okB {
		return tA4.Value.Lexeme == tB4.Value.Lexeme
	}

	tA5, okA := a.(ast.StringLiteralExpr)
	tB5, okB := b.(ast.StringLiteralExpr)

	if okA && okB {
		return tA5.Value.Lexeme == tB5.Value.Lexeme
	}

	_, okA = a.(ast.BottomLiteralExpr)
	_, okB = b.(ast.BottomLiteralExpr)

	if okA && okB {
		return true
	}

	tA6, okA := a.(ast.IdentifierExpr)
	tB6, okB := b.(ast.IdentifierExpr)

	if okA && okB {
		return tA6.Name.Lexeme == tB6.Name.Lexeme
	}

	tA7, okA := a.(ast.VarDeclStmt)
	tB7, okB := b.(ast.VarDeclStmt)

	if okA && okB {
		if len(tA7.Names) != len(tB7.Names) {
			return false
		}

		if len(tA7.Values) != len(tB7.Values) {
			return false
		}

		for i, n1 := range tA7.Names {
			n2 := tB7.Names[i]

			if n1.Name.Lexeme != n2.Name.Lexeme || n1.Mut != n2.Mut {
				return false
			}
		}

		for i, v1 := range tA7.Values {
			v2 := tB7.Values[i]

			if !nodesAreEqual(v1, v2) {
				return false
			}
		}

		return true
	}

	tA8, okA := a.(ast.BooleanLiteralExpr)
	tB8, okB := b.(ast.BooleanLiteralExpr)

	if okA && okB {
		return tA8.Value.Type == tB8.Value.Type
	}

	tA9, okA := a.(ast.AssignmentStmt)
	tB9, okB := b.(ast.AssignmentStmt)

	if okA && okB {
		for i, n := range tA9.Names {
			if n.Lexeme != tB9.Names[i].Lexeme {
				return false
			}
		}

		for i, n := range tA9.Values {
			if !nodesAreEqual(n, tB9.Values[i]) {
				return false
			}
		}

		return true
	}

	tA10, okA := a.(ast.IfStmt)
	tB10, okB := b.(ast.IfStmt)

	if okA && okB {
		return nodesAreEqual(tA10.Decls, tB10.Decls) &&
			nodesAreEqual(tA10.Condition, tB10.Condition) &&
			nodesAreEqual(tA10.Then, tB10.Then) &&
			nodesAreEqual(tA10.Else, tB10.Else)
	}

	tA11, okA := a.(ast.Block)
	tB11, okB := b.(ast.Block)

	if okA && okB {
		if len(tA11.Contents) != len(tB11.Contents) {
			return false
		}

		for i, n := range tA11.Contents {
			if !nodesAreEqual(n, tB11.Contents[i]) {
				return false
			}
		}

		return true
	}

	tA12, okA := a.(ast.FuncExpr)
	tB12, okB := b.(ast.FuncExpr)

	if okA && okB {
		for i, v := range tA12.Params {
			if v.Name.Lexeme != tB12.Params[i].Name.Lexeme || v.Mut != tB12.Params[i].Mut {
				return false
			}
		}

		return nodesAreEqual(tA12.Body, tB12.Body)
	}

	tA13, okA := a.(ast.ReturnStmt)
	tB13, okB := b.(ast.ReturnStmt)

	if okA && okB {
		return nodesAreEqual(tA13.Expr, tB13.Expr)
	}

	tA14, okA := a.(ast.CallExpr)
	tB14, okB := b.(ast.CallExpr)

	if okA && okB {
		if len(tA14.Arguments) != len(tB14.Arguments) {
			return false
		}

		for i, a := range tA14.Arguments {
			b := tB14.Arguments[i]

			if !nodesAreEqual(a, b) {
				return false
			}
		}

		if !nodesAreEqual(tA14.Callee, tB14.Callee) {
			return false
		}

		return true
	}

	tA15, okA := a.(ast.TupleLiteralExpr)
	tB15, okB := b.(ast.TupleLiteralExpr)

	if okA && okB {
		if len(tA15.Contents) != len(tB15.Contents) {
			return false
		}

		for i, e := range tA15.Contents {
			if !nodesAreEqual(e, tB15.Contents[i]) {
				return false
			}
		}

		return true
	}

	_, okA = a.(ast.MeExpr)
	_, okB = b.(ast.MeExpr)

	if okA && okB {
		return true
	}

	tA16, okA := a.(ast.ProtoExpr)
	tB16, okB := b.(ast.ProtoExpr)

	if okA && okB {
		if len(tA16.MethodSet) != len(tB16.MethodSet) {
			return false
		}

		for i, m := range tA16.MethodSet {
			if !nodesAreEqual(m.K, tB16.MethodSet[i].K) ||
				!nodesAreEqual(m.M, tB16.MethodSet[i].M) {
				return false
			}
		}

		return true
	}

	tA17, okA := a.(ast.GetExpr)
	tB17, okB := b.(ast.GetExpr)

	if okA && okB {
		return nodesAreEqual(tA17.Gettee, tB17.Gettee) &&
			nodesAreEqual(tA17.Field, tB17.Field)
	}

	return false
}

func astsAreEqual(as []ast.Node, bs []ast.Node) bool {
	if len(as) != len(bs) {
		return false
	}

	for i, a := range as {
		b := bs[i]

		if !nodesAreEqual(a, b) {
			return false
		}
	}

	return true
}

func TestParse(t *testing.T) {
	t.Run("productions", func(t *testing.T) {
		table := []struct {
			name     string
			text     string
			expected []ast.Node
		}{
			{
				name:     "fundamental string 1",
				text:     "'abc'",
				expected: []ast.Node{ast.StringLiteralExpr{Value: tokens.New(tokentype.STRING, "'abc'", 0, 0)}},
			},
			{
				name:     "fundamental string 2",
				text:     "\"abc\"",
				expected: []ast.Node{ast.StringLiteralExpr{Value: tokens.New(tokentype.STRING, "\"abc\"", 0, 0)}},
			},
			{
				name:     "fundamental number",
				text:     "123",
				expected: []ast.Node{ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "123", 0, 0)}},
			},
			{
				name:     "fundamental grouping",
				text:     "(123)",
				expected: []ast.Node{ast.GroupingExpr{Expr: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "123", 0, 0)}}},
			},
			{
				name:     "fundamental identifier",
				text:     "abc",
				expected: []ast.Node{ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "abc", 0, 0)}},
			},
			{
				name:     "fundamental bottom",
				text:     "bottom",
				expected: []ast.Node{ast.BottomLiteralExpr{Token: tokens.New(tokentype.BOTTOM, "bottom", 0, 0)}},
			},
			{
				name:     "fundamental boolean 1",
				text:     "true",
				expected: []ast.Node{ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)}},
			},
			{
				name:     "fundamental boolean 2",
				text:     "false",
				expected: []ast.Node{ast.BooleanLiteralExpr{Value: tokens.New(tokentype.FALSE, "false", 0, 0)}},
			},
			{
				name: "fundamental function 1",
				text: "fn (x, mut y) -> true",
				expected: []ast.Node{
					ast.FuncExpr{
						Params: []ast.Identifier{
							{Name: tokens.New(tokentype.IDENTIFIER, "x", 0, 0), Mut: false},
							{Name: tokens.New(tokentype.IDENTIFIER, "y", 0, 0), Mut: true},
						},
						Body: ast.Block{
							Contents: []ast.Node{
								ast.ReturnStmt{Expr: ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)}},
							},
						},
					},
				},
			},
			{
				name: "fundamental function 2",
				text: "fn (a, mut b) { true }",
				expected: []ast.Node{
					ast.FuncExpr{
						Params: []ast.Identifier{
							{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0), Mut: false},
							{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 0), Mut: true},
						},
						Body: ast.Block{
							Contents: []ast.Node{
								ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
							},
						},
					},
				},
			},
			{
				name: "fundamental tuple 1",
				text: "[1]",
				expected: []ast.Node{
					ast.TupleLiteralExpr{
						Contents: []ast.Expr{
							ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						},
					},
				},
			},
			{
				name: "fundamental tuple 2",
				text: "[]",
				expected: []ast.Node{
					ast.TupleLiteralExpr{
						Contents: nil,
					},
				},
			},
			{
				name: "fundamental tuple 3",
				text: "[1+2, a, 'a']",
				expected: []ast.Node{
					ast.TupleLiteralExpr{
						Contents: []ast.Expr{
							ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
							ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
							ast.StringLiteralExpr{Value: tokens.New(tokentype.STRING, "'a'", 0, 0)},
						},
					},
				},
			},
			{
				name: "funamental me",
				text: "me",
				expected: []ast.Node{
					ast.MeExpr{Token: tokens.New(tokentype.ME, "me", 0, 0)},
				},
			},
			{
				name: "proto 1",
				text: "proto { 4 -> fn () -> 1 }",
				expected: []ast.Node{
					ast.ProtoExpr{
						MethodSet: []ast.ProtoMethod{
							{
								K: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
								M: ast.FuncExpr{Body: ast.Block{
									Contents: []ast.Node{
										ast.ReturnStmt{Expr: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)}},
									},
								}},
							},
						},
					},
				},
			},
			{
				name: "proto 2",
				text: "proto { 4 -> fn () -> 1, true -> fn () { return 3; } }",
				expected: []ast.Node{
					ast.ProtoExpr{
						MethodSet: []ast.ProtoMethod{
							{
								K: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
								M: ast.FuncExpr{Body: ast.Block{
									Contents: []ast.Node{
										ast.ReturnStmt{Expr: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)}},
									},
								}},
							},
							{
								K: ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
								M: ast.FuncExpr{Body: ast.Block{
									Contents: []ast.Node{
										ast.ReturnStmt{Expr: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)}},
									},
								}},
							},
						},
					},
				},
			},
			{
				name: "call expression 1",
				text: "fn () {}()",
				expected: []ast.Node{
					ast.CallExpr{
						Callee: ast.FuncExpr{
							Params: []ast.Identifier{},
							Body: ast.Block{
								Contents: []ast.Node{},
							},
						},
					},
				},
			},
			{
				name: "call expression 2",
				text: "(fn () -> 1)()",
				expected: []ast.Node{
					ast.CallExpr{
						Callee: ast.GroupingExpr{
							Expr: ast.FuncExpr{
								Body: ast.Block{
									Contents: []ast.Node{
										ast.ReturnStmt{
											Expr: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
										},
									},
								},
							},
						},
						Arguments: []ast.Expr{},
					},
				},
			},
			{
				name: "call expression 3",
				text: "abc()",
				expected: []ast.Node{
					ast.CallExpr{
						Callee:    ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "abc", 0, 0)},
						Arguments: []ast.Expr{},
					},
				},
			},
			{
				name: "call expression 4",
				text: "abc()()",
				expected: []ast.Node{
					ast.CallExpr{
						Callee: ast.CallExpr{
							Callee:    ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "abc", 0, 0)},
							Arguments: []ast.Expr{},
						},
						Arguments: []ast.Expr{},
					},
				},
			},
			{
				name: "call expression 5",
				text: "a(true, 1 + 2)",
				expected: []ast.Node{
					ast.CallExpr{
						Callee: ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
						Arguments: []ast.Expr{
							ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
							ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
						},
					},
				},
			},
			{
				name: "get expression 1",
				text: "abc->def",
				expected: []ast.Node{
					ast.GetExpr{
						Gettee: ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "abc", 0, 0)},
						Field:  ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "def", 0, 0)},
					},
				},
			},
			{
				name: "get expression 2",
				text: "[]->'abc'->3->true",
				expected: []ast.Node{
					ast.GetExpr{
						Gettee: ast.GetExpr{
							Gettee: ast.GetExpr{
								Gettee: ast.TupleLiteralExpr{},
								Field:  ast.StringLiteralExpr{Value: tokens.New(tokentype.STRING, "'abc'", 0, 0)},
							},
							Field: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						},
						Field: ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
					},
				},
			},
			{
				name: "get expression 3",
				text: "[]->'abc'()->'def'()",
				expected: []ast.Node{
					ast.CallExpr{
						Callee: ast.GetExpr{
							Gettee: ast.CallExpr{
								Callee: ast.GetExpr{
									Gettee: ast.TupleLiteralExpr{Contents: nil},
									Field:  ast.StringLiteralExpr{Value: tokens.New(tokentype.STRING, "'abc'", 0, 0)},
								},
								Arguments: []ast.Expr{},
							},
							Field: ast.StringLiteralExpr{Value: tokens.New(tokentype.STRING, "'def'", 0, 0)},
						},
						Arguments: []ast.Expr{},
					},
				},
			},
			{
				name: "unary minus",
				text: "-7",
				expected: []ast.Node{
					ast.UnaryExpr{
						Operator: tokens.New(tokentype.MINUS, "-", 0, 0),
						Expr:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "7", 0, 0)},
					},
				},
			},
			{
				name: "binary exponentiation",
				text: "1 ** 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.ASTERISK_ASTERISK, "**", 0, 0),
					},
				},
			},
			{
				name: "binary multiplication",
				text: "1 * 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.ASTERISK, "*", 0, 0),
					},
				},
			},
			{
				name: "binary division",
				text: "1 / 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.SLASH, "/", 0, 0),
					},
				},
			},
			{
				name: "binary addition",
				text: "1 + 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
					},
				},
			},
			{
				name: "binary subtraction",
				text: "1 - 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.MINUS, "-", 0, 0),
					},
				},
			},
			{
				name: "binary comparison 1",
				text: "4 < 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.LESS, "<", 0, 0),
					},
				},
			},
			{
				name: "binary comparison 2",
				text: "4 <= 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.LESS_EQUAL, "<=", 0, 0),
					},
				},
			},
			{
				name: "binary comparison 3",
				text: "4 > 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.GREAT, ">", 0, 0),
					},
				},
			},
			{
				name: "binary comparison 4",
				text: "4 >= 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.GREAT_EQUAL, ">=", 0, 0),
					},
				},
			},
			{
				name: "equality comparison 1",
				text: "1 == 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
					},
				},
			},
			{
				name: "equality comparison 2",
				text: "1 != 2",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						Operator: tokens.New(tokentype.BANG_EQUAL, "!=", 0, 0),
					},
				},
			},
			{
				name: "boolean and",
				text: "true && false",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
						Right:    ast.BooleanLiteralExpr{Value: tokens.New(tokentype.FALSE, "false", 0, 0)},
						Operator: tokens.New(tokentype.AMPERSAND_AMPERSAND, "&&", 0, 0),
					},
				},
			},
			{
				name: "boolean or",
				text: "true || false",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left:     ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
						Right:    ast.BooleanLiteralExpr{Value: tokens.New(tokentype.FALSE, "false", 0, 0)},
						Operator: tokens.New(tokentype.STROKE_STROKE, "||", 0, 0),
					},
				},
			},
			{
				name: "variable declaration w/o inits",
				text: "let a, b;",
				expected: []ast.Node{
					ast.VarDeclStmt{
						Names: []ast.Identifier{
							{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0), Mut: false},
							{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 0), Mut: false},
						},
						Values: []ast.Expr{},
					},
				},
			},
			{
				name: "variable declaration w/ inits",
				text: "let a, b = 1, 2;",
				expected: []ast.Node{
					ast.VarDeclStmt{
						Names: []ast.Identifier{
							{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0), Mut: false},
							{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 0), Mut: false},
						},
						Values: []ast.Expr{
							ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						},
					},
				},
			},
			{
				name: "mutable variable declaration",
				text: "let mut a;",
				expected: []ast.Node{
					ast.VarDeclStmt{
						Names: []ast.Identifier{
							{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0), Mut: true},
						},
						Values: []ast.Expr{},
					},
				},
			},
			{
				name: "assignment expression w/ one entry",
				text: "a = 1;",
				expected: []ast.Node{
					ast.AssignmentStmt{
						Names:  []tokens.Token{tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
						Values: []ast.Expr{ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)}},
					},
				},
			},
			{
				name: "assignment expression w/ multiple entries",
				text: "a, b = 1, 2 + 5;",
				expected: []ast.Node{
					ast.AssignmentStmt{
						Names: []tokens.Token{
							tokens.New(tokentype.IDENTIFIER, "a", 0, 0),
							tokens.New(tokentype.IDENTIFIER, "b", 0, 0),
						},
						Values: []ast.Expr{
							ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "5", 0, 0)},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
						},
					},
				},
			},
			{
				name: "if with just condition",
				text: "if a == 4 {}",
				expected: []ast.Node{
					ast.IfStmt{
						Decls: ast.VarDeclStmt{},
						Condition: ast.BinaryExpr{
							Left:     ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
							Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
						},
						Then: ast.Block{Contents: []ast.Node{}},
						Else: nil,
					},
				},
			},
			{
				name: "if with declaration and condition",
				text: "if let a = 1; a == 4 {}",
				expected: []ast.Node{
					ast.IfStmt{
						Decls: ast.VarDeclStmt{
							Names: []ast.Identifier{
								{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0), Mut: false},
							},
							Values: []ast.Expr{ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)}},
						},
						Condition: ast.BinaryExpr{
							Left:     ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
							Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
						},
						Then: ast.Block{Contents: []ast.Node{}},
						Else: nil,
					},
				},
			},
			{
				name: "if with condition and then body",
				text: "if a == 4 { 1 + 1 }",
				expected: []ast.Node{
					ast.IfStmt{
						Decls: ast.VarDeclStmt{},
						Condition: ast.BinaryExpr{
							Left:     ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
							Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
						},
						Then: ast.Block{
							Contents: []ast.Node{
								ast.BinaryExpr{
									Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
									Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
									Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
								},
							},
						},
						Else: nil,
					},
				},
			},
			{
				name: "if with condition and else body",
				text: "if a == 4 {} else { 1 + 1 }",
				expected: []ast.Node{
					ast.IfStmt{
						Decls: ast.VarDeclStmt{},
						Condition: ast.BinaryExpr{
							Left:     ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
							Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
						},
						Then: ast.Block{Contents: []ast.Node{}},
						Else: ast.Block{Contents: []ast.Node{
							ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
						}},
					},
				},
			},
			{
				name: "if with condition and else if",
				text: "if a == 4 {} else if b == 2 {}",
				expected: []ast.Node{
					ast.IfStmt{
						Decls: ast.VarDeclStmt{},
						Condition: ast.BinaryExpr{
							Left:     ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
							Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
						},
						Then: ast.Block{Contents: []ast.Node{}},
						Else: ast.IfStmt{
							Decls: ast.VarDeclStmt{},
							Condition: ast.BinaryExpr{
								Left:     ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "b", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
								Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
							},
							Then: ast.Block{Contents: []ast.Node{}},
							Else: nil,
						},
					},
				},
			},
			{
				name: "return 1",
				text: "return;",
				expected: []ast.Node{
					ast.ReturnStmt{Expr: nil},
				},
			},
			{
				name: "return 2",
				text: "return false;",
				expected: []ast.Node{
					ast.ReturnStmt{Expr: ast.BooleanLiteralExpr{Value: tokens.New(tokentype.FALSE, "false", 0, 0)}},
				},
			},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: got error unexpectedly during scanning", e.name)
			}

			ast, err := parser.New(ts).Parse()

			if err != nil {
				t.Errorf("%q: received unexpected parse error", e.name)
			}

			if !astsAreEqual(ast, e.expected) {
				t.Errorf("%q: generated ast %#v does not equal expected %#v", e.name, ast, e.expected)
			}
		}
	})

	t.Run("operator associativity", func(t *testing.T) {
		table := []struct {
			name     string
			text     string
			expected []ast.Node
		}{
			{
				name: "left associativity addition",
				text: "1 + 2 - 3 + 4",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
							Operator: tokens.New(tokentype.MINUS, "-", 0, 0),
						},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
						Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
					},
				},
			},
			{
				name: "left associativity multiplication",
				text: "1 * 2 / 3 * 4",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
								Operator: tokens.New(tokentype.ASTERISK, "*", 0, 0),
							},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
							Operator: tokens.New(tokentype.SLASH, "/", 0, 0),
						},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
						Operator: tokens.New(tokentype.ASTERISK, "*", 0, 0),
					},
				},
			},
			{
				name: "left associativity exponentiation",
				text: "1 ** 2 ** 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
							Operator: tokens.New(tokentype.ASTERISK_ASTERISK, "**", 0, 0),
						},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.ASTERISK_ASTERISK, "**", 0, 0),
					},
				},
			},
			{
				name: "right associativity unary",
				text: "--1",
				expected: []ast.Node{
					ast.UnaryExpr{
						Expr: ast.UnaryExpr{
							Expr:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							Operator: tokens.New(tokentype.MINUS, "-", 0, 0),
						},
						Operator: tokens.New(tokentype.MINUS, "-", 0, 0),
					},
				},
			},
			{
				name: "left associativity boolean and",
				text: "1 && 2 && 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
							Operator: tokens.New(tokentype.AMPERSAND_AMPERSAND, "&&", 0, 0),
						},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.AMPERSAND_AMPERSAND, "&&", 0, 0),
					},
				},
			},
			{
				name: "left associativity boolean or",
				text: "1 || 2 || 3",
				expected: []ast.Node{
					ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
							Operator: tokens.New(tokentype.STROKE_STROKE, "||", 0, 0),
						},
						Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
						Operator: tokens.New(tokentype.STROKE_STROKE, "||", 0, 0),
					},
				},
			},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: got error unexpectedly during scanning", e.name)
			}

			ast, err := parser.New(ts).Parse()

			if err != nil {
				t.Errorf("%q: received unexpected parse error", e.name)
			}

			if !astsAreEqual(ast, e.expected) {
				t.Errorf("%q: generated ast %#v does not equal expected %#v", e.name, ast, e.expected)
			}
		}
	})

	t.Run("precedence", func(t *testing.T) {
		expected := []ast.Node{
			ast.BinaryExpr{
				Left: ast.BooleanLiteralExpr{Value: tokens.New(tokentype.TRUE, "true", 0, 0)},
				Right: ast.BinaryExpr{
					Left: ast.BooleanLiteralExpr{Value: tokens.New(tokentype.FALSE, "false", 0, 0)},
					Right: ast.BinaryExpr{
						Left: ast.IdentifierExpr{Name: tokens.New(tokentype.IDENTIFIER, "a", 0, 0)},
						Right: ast.BinaryExpr{
							Left: ast.BinaryExpr{
								Left: ast.UnaryExpr{
									Operator: tokens.New(tokentype.MINUS, "-", 0, 0),
									Expr:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								},
								Right: ast.BinaryExpr{
									Left: ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
									Right: ast.BinaryExpr{
										Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "3", 0, 0)},
										Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "4", 0, 0)},
										Operator: tokens.New(tokentype.ASTERISK_ASTERISK, "**", 0, 0),
									},
									Operator: tokens.New(tokentype.ASTERISK, "*", 0, 0),
								},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
							Right: ast.BinaryExpr{
								Left:     ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
								Right:    ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "5", 0, 0)},
								Operator: tokens.New(tokentype.PLUS, "+", 0, 0),
							},
							Operator: tokens.New(tokentype.LESS, "<", 0, 0),
						},
						Operator: tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0),
					},
					Operator: tokens.New(tokentype.AMPERSAND_AMPERSAND, "&&", 0, 0),
				},
				Operator: tokens.New(tokentype.STROKE_STROKE, "||", 0, 0),
			},
		}

		ts, err := scanner.New().Read("true || false && a == -1 + 2 * 3 ** 4 < 1 + 5")

		if err != nil {
			t.Error("got error unexpectedly during scanning")
		}

		ast, err := parser.New(ts).Parse()

		if err != nil {
			t.Error("received unexpected parse error")
		}

		if !astsAreEqual(ast, expected) {
			t.Errorf("generated ast %#v does not equal expected %#v", ast, expected)
		}
	})

	t.Run("parse errors", func(t *testing.T) {
		table := []struct {
			name string
			text string
		}{
			{name: "malformed grouping expression 1", text: "(1 + 2"},
			{name: "malformed grouping expression 2", text: "(1 + 2}"},
			{name: "malformed expression in grouping", text: "(1 +)"},
			{name: "malformed unary expression", text: "-+"},
			{name: "malformed exponentiation expression 1", text: "1 **"},
			{name: "malformed exponentiation expression 2", text: "1 ** -"},
			{name: "malformed multiplication expression 1", text: "1 *"},
			{name: "malformed multiplication expression 2", text: "1 * *"},
			{name: "malformed addition expression 1", text: "1 +"},
			{name: "malformed addition expression 2", text: "1 + -"},
			{name: "malformed comparison expression", text: "1 < *"},
			{name: "malformed equality expression", text: "1 == +"},
			{name: "malformed boolean and expression", text: "true && !"},
			{name: "malformed boolean or expression", text: "true || !"},
			{name: "malformed function expression 1", text: "fn a -> 1"},
			{name: "malformed function expression 2", text: "fn (1) -> 1"},
			{name: "malformed function expression 3", text: "fn (a -> 1"},
			{name: "malformed function expression 3", text: "fn (a) -> 1 +"},
			{name: "malformed function expression 3", text: "fn (a) { 1 + }"},
			{name: "malformed call expression", text: "a(if)"},
			{name: "malformed get expression", text: "1->if true {}"},
			{name: "malformed variable declaration 1", text: "let 'ab';"},
			{name: "malformed variable declaration 2", text: "let ab"},
			{name: "malformed variable declaration 3", text: "let a = 4 +"},
			{name: "malformed variable declaration 4", text: "let a = 'b'"},
			{name: "malformed variable declaration 5", text: "let a, 1 = 'b'"},
			{name: "malformed variable declaration 6", text: "let a, b = 1, ;"},
			{name: "malformed assignment statemement 1", text: "a = 1 +"},
			{name: "malformed assignment statemement 2", text: "1 = 1"},
			{name: "malformed assignment statemement 3", text: "a, 2 = 1 +"},
			{name: "malformed assignment statemement 4", text: "a, b"},
			{name: "malformed assignment statemement 5", text: "a = 1 + 2"},
			{name: "malformed if statment variable declaration", text: "if let; true {}"},
			{name: "malformed if statment condition", text: "if 1 + {}"},
			{name: "malformed if statment `then` block", text: "if true {"},
			{name: "malformed if statment `else` block 1", text: "if true {} else {"},
			{name: "malformed if statment `else` block 2", text: "if true {} else }"},
			{name: "malformed if statment `else` block 2", text: "if true {} else }"},
			{name: "malformed if statment `else` if block", text: "if true {} else if {}"},
			{name: "malformed return statement", text: "return true"},
			{name: "malformed return statement", text: "return 1 +;"},
			{name: "malformed tuple expression 1", text: "[1,]"},
			{name: "malformed tuple expression 2", text: "[1+]"},
			{name: "malformed tuple expression 3", text: "[1"},
			{name: "malformed proto expression 1", text: "proto { }"},
			{name: "malformed proto expression 2", text: "proto 'a' -> fn () -> 1 }"},
			{name: "malformed proto expression 3", text: "proto { 'a' -> fn () -> 1"},
			{name: "malformed proto expression 4", text: "proto { 'a' -> 2 }"},
			{name: "malformed proto expression 5", text: "proto { 'a' -> fn () -> 1, 'b' -> 3 }"},
			{name: "malformed proto expression 6", text: "proto { 'a' fn () -> 1 }"},
			{name: "malformed proto expression 7", text: "proto { 'a' -> fn -> 1 }"},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: got error unexpectedly during scanning", e.name)
			}

			_, err = parser.New(ts).Parse()

			if err == nil {
				t.Errorf("%q: did not receive an error on malformed input %q", e.name, e.text)
			}
		}
	})
}
