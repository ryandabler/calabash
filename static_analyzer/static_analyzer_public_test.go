package staticanalyzer_test

import (
	"calabash/lexer/scanner"
	"calabash/parser"
	staticanalyzer "calabash/static_analyzer"
	"testing"
)

func TestAnalyze(t *testing.T) {
	t.Run("corrects", func(t *testing.T) {
		table := []struct {
			name string
			text string
		}{
			{
				name: "numeric literals",
				text: "1",
			},
			{
				name: "string literals",
				text: "\"a\"",
			},
			{
				name: "bottom literals",
				text: "bottom",
			},
			{
				name: "unary expressions",
				text: "-1",
			},
			{
				name: "grouping expressions",
				text: "('abc')",
			},
			{
				name: "binary expressions",
				text: "1 + 3",
			},
			{
				name: "boolean expressions 1",
				text: "true",
			},
			{
				name: "boolean expressions 2",
				text: "false",
			},
			{
				name: "variable declaration",
				text: "let a;",
			},
			{
				name: "identifier expressions",
				text: "let a; a",
			},
			{
				name: "assignment statement",
				text: "let a; a = 1;",
			},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: got unexpected lexer error %q", e.name, err.Error())
			}

			ast, err := parser.New(ts).Parse()

			if err != nil {
				t.Errorf("%q: got unexpected parse error %q", e.name, err.Error())
			}

			err = staticanalyzer.New().Analyze(ast)

			if err != nil {
				t.Errorf("%q: go unexpected static error %q", e.name, err.Error())
			}
		}
	})

	t.Run("errors", func(t *testing.T) {
		table := []struct {
			name string
			text string
		}{
			{
				name: "undeclared identifier",
				text: "a",
			},
			{
				name: "under-initialized variables",
				text: "let a, b = 1;",
			},
			{
				name: "redeclaring variable",
				text: "let a; let a;",
			},
			{
				name: "left-side binary error",
				text: "a + 2",
			},
			{
				name: "right-side binary error",
				text: "2 - a",
			},
			{
				name: "var declaration with undeclared identifier expression",
				text: "let a = b;",
			},
			{
				name: "assignment name/value quantity mismatch",
				text: "let a, b; a, b = 1;",
			},
			{
				name: "assignment with undeclared variable",
				text: "a = 1;",
			},
			{
				name: "assignment with undeclared identifier expression",
				text: "let a; a = b;",
			},
		}

		for _, e := range table {
			ts, err := scanner.New().Read(e.text)

			if err != nil {
				t.Errorf("%q: got unexpected lexer error %q", e.name, err.Error())
			}

			ast, err := parser.New(ts).Parse()

			if err != nil {
				t.Errorf("%q: got unexpected parse error %q", e.name, err.Error())
			}

			err = staticanalyzer.New().Analyze(ast)

			if err == nil {
				t.Errorf("%q: did not receive static analysis error for program %q", e.name, e.text)
			}
		}
	})
}
