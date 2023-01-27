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

			if n1.Lexeme != n2.Lexeme {
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
				name: "variable declaration w/o inits",
				text: "let a, b;",
				expected: []ast.Node{
					ast.VarDeclStmt{
						Names:  []tokens.Token{tokens.New(tokentype.IDENTIFIER, "a", 0, 0), tokens.New(tokentype.IDENTIFIER, "b", 0, 0)},
						Values: []ast.Expr{},
					},
				},
			},
			{
				name: "variable declaration w/ inits",
				text: "let a, b = 1, 2;",
				expected: []ast.Node{
					ast.VarDeclStmt{
						Names: []tokens.Token{tokens.New(tokentype.IDENTIFIER, "a", 0, 0), tokens.New(tokentype.IDENTIFIER, "b", 0, 0)},
						Values: []ast.Expr{
							ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "1", 0, 0)},
							ast.NumericLiteralExpr{Value: tokens.New(tokentype.NUMBER, "2", 0, 0)},
						},
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
		}

		ts, err := scanner.New().Read("-1 + 2 * 3 ** 4")

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
			{name: "malformed variable declaration 1", text: "let 'ab';"},
			{name: "malformed variable declaration 2", text: "let ab"},
			{name: "malformed variable declaration 3", text: "let a = 4 +"},
			{name: "malformed variable declaration 4", text: "let a = 'b'"},
			{name: "malformed variable declaration 5", text: "let a, 1 = 'b'"},
			{name: "malformed variable declaration 6", text: "let a, b = 1, ;"},
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
