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
				name: "tuple expressions 1",
				text: "[]",
			},
			{
				name: "tuple expressions 2",
				text: "[1, 'a', fn () -> bottom]",
			},
			{
				name: "tuple expressions 3",
				text: "let a; [a]",
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
				name: "function expressions",
				text: "fn (a, mut b) -> a + b",
			},
			{
				name: "call expressions 1",
				text: "let abc; abc()",
			},
			{
				name: "call expressions 2",
				text: "let a, abc; abc(a)",
			},
			{
				name: "proto expression",
				text: "proto { true -> fn(a) -> a }",
			},
			{
				name: "me expression",
				text: "proto { true -> fn() -> me }",
			},
			{
				name: "call expression 3",
				text: "fn () {}()",
			},
			{
				name: "record expressions 1",
				text: "{ 1 -> 1, 'a' -> 1, true -> 1, bottom -> 1, [1] -> 1, {} -> 1, proto { 1 -> fn () {} } -> 1, ('a' + 'b') -> 1 }",
			},
			{
				name: "record expressions 2",
				text: "let a = 1; { a -> a }",
			},
			{
				name: "get expression 1",
				text: "1->'b'",
			},
			{
				name: "pipe expression 1",
				text: "1 |> 2 + ?",
			},
			{
				name: "pipe expression 2",
				text: "1 |> (2 |> ?) + ?",
			},
			{
				name: "assignment statement",
				text: "let mut a; a = 1;",
			},
			{
				name: "if statement with only condition",
				text: "let a; if a == 4 {}",
			},
			{
				name: "if statement with declaration",
				text: "if let a; a == 4 {}",
			},
			{
				name: "if statement with shadowing declaration",
				text: "let a; if let a; a == 4 {}",
			},
			{
				name: "if statement with declaration and shadowing declaration in `then` block",
				text: "if let a; a == 4 { let a; }",
			},
			{
				name: "if statement with else...if block",
				text: "let a; if a == 4 {} else if a == 3 {}",
			},
			{
				name: "if statement with shadowing else...if block",
				text: "if let a; a == 4 {} else if let a; a == 3 {}",
			},
			{
				name: "while statement with no variable declaration",
				text: "while true {}",
			},
			{
				name: "while statement with variable declaration",
				text: "while let a = 1; a == 1 {}",
			},
			{
				name: "while statement with shadowing variable declaration",
				text: "let a = 1; while let a = 1; true {}",
			},
			{
				name: "while statement with variable lookup",
				text: "let a = 1; while a == 1 {}",
			},
			{
				name: "while statement with variable lookup in block",
				text: "let a = 1; while a == 1 { a }",
			},
			{
				name: "return statement in function",
				text: "fn () { return 1; }",
			},
			{
				name: "return statement in proto methods",
				text: "proto { 'a' -> fn() { return 2; } }",
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
				name: "referencing identifier not in immediate scope",
				text: "let a; fn (b) -> a",
			},
			{
				name: "referencing identifier not declared in body",
				text: "let a; fn (b) { a }",
			},
			{
				name: "referencing undeclared identifier in arguments list",
				text: "fn () {}(a)",
			},
			{
				name: "calling an undeclared identifier",
				text: "abc()",
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
				name: "me used outside of proto",
				text: "me",
			},
			{
				name: "proto with 'me' used as key",
				text: "proto { me -> fn () -> 1 }",
			},
			{
				name: "proto with undeclared identifier as key",
				text: "proto { a -> fn() -> 1 }",
			},
			{
				name: "proto with undeclared identifier in method",
				text: "proto { 'a' -> fn() -> a }",
			},
			{
				name: "record with undeclared identifier key",
				text: " { a -> 1 }",
			},
			{
				name: "record with undeclared identifier value",
				text: " { 1 -> a }",
			},
			{
				name: "get expression with unresolved gettee variable",
				text: "a -> 1",
			},
			{
				name: "get expression with unresolved field variable",
				text: "'a' -> a",
			},
			{
				name: "question mark outside of pipe expression",
				text: "1 + ?",
			},
			{
				name: "question mark not referenced in a pipe expression",
				text: "1 |> 1",
			},
			{
				name: "question mark not referenced in an outer pipe expression",
				text: "1 |> (1 |> ?)",
			},
			{
				name: "question mark not referenced in an inner pipe expression",
				text: "1 |> ? + (1 |> 1)",
			},
			{
				name: "var declaration with undeclared identifier expression",
				text: "let a = b;",
			},
			{
				name: "assignment name/value quantity mismatch",
				text: "let mut a, mut b; a, b = 1;",
			},
			{
				name: "assignment with undeclared variable",
				text: "a = 1;",
			},
			{
				name: "assignment to immutable variable",
				text: "let a; a = 1;",
			},
			{
				name: "assignment with undeclared identifier expression",
				text: "let mut a; a = b;",
			},
			{
				name: "if statement with undeclared identifier in variable declaration",
				text: "if let a = b; a == 3 {}",
			},
			{
				name: "if statement with undeclared identifier in condition",
				text: "if a == 3 {}",
			},
			{
				name: "if statement with multiple redeclarations in `then` block",
				text: "if let a; a == 3 { let a, a; }",
			},
			{
				name: "if statement with undeclared identifier in `else` block",
				text: "if let a; a == 3 {} else { b }",
			},
			{
				name: "tuple expression containing undeclared variables",
				text: "[1, a]",
			},
			{
				name: "while statement with unresolved condition variable",
				text: "while a == 1 {}",
			},
			{
				name: "while statement with unresolved block variable",
				text: "while true { a }",
			},
			{
				name: "top level return statements",
				text: "return 3;",
			},
			{
				name: "return statements in loop",
				text: "while true { return 3; }",
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
