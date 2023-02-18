package scanner_test

import (
	"calabash/internal/tokentype"
	"calabash/lexer/scanner"
	"calabash/lexer/tokens"
	"testing"
)

func same(as []tokens.Token, bs []tokens.Token) bool {
	if len(as) != len(bs) {
		return false
	}

	if len(as) == 0 && len(bs) == 0 {
		return true
	}

	for i, a := range as {
		b := bs[i]

		if a.Type != b.Type || a.Lexeme != b.Lexeme {
			return false
		}
	}

	return true
}

func samePos(as []tokens.Token, bs []tokens.Token) bool {
	for i, a := range as {
		b := bs[i]

		if a.Position.Col != b.Position.Col || a.Position.Row != b.Position.Row {
			return false
		}
	}

	return true
}

func lessEOF(ts []tokens.Token) []tokens.Token {
	if len(ts) == 0 {
		return ts
	}

	if ts[len(ts)-1].Type == tokentype.EOF {
		return ts[:len(ts)-1]
	}

	return ts
}

func TestRead(t *testing.T) {
	/*
		Produce basic tokens
	*/
	table := []struct {
		name      string
		text      string
		expected  []tokens.Token
		willError bool
	}{
		{name: "left paren", text: "(", expected: []tokens.Token{tokens.New(tokentype.LEFT_PAREN, "(", 0, 0)}},
		{name: "right paren", text: ")", expected: []tokens.Token{tokens.New(tokentype.RIGHT_PAREN, ")", 0, 0)}},
		{name: "left bracket", text: "[", expected: []tokens.Token{tokens.New(tokentype.LEFT_BRACKET, "[", 0, 0)}},
		{name: "right bracket", text: "]", expected: []tokens.Token{tokens.New(tokentype.RIGHT_BRACKET, "]", 0, 0)}},
		{name: "left brace", text: "{", expected: []tokens.Token{tokens.New(tokentype.LEFT_BRACE, "{", 0, 0)}},
		{name: "right brace", text: "}", expected: []tokens.Token{tokens.New(tokentype.RIGHT_BRACE, "}", 0, 0)}},
		{name: "comma", text: ",", expected: []tokens.Token{tokens.New(tokentype.COMMA, ",", 0, 0)}},
		{name: "semicolon", text: ";", expected: []tokens.Token{tokens.New(tokentype.SEMICOLON, ";", 0, 0)}},
		{name: "less", text: "<", expected: []tokens.Token{tokens.New(tokentype.LESS, "<", 0, 0)}},
		{name: "less equal", text: "<=", expected: []tokens.Token{tokens.New(tokentype.LESS_EQUAL, "<=", 0, 0)}},
		{name: "less less", text: "<<", expected: []tokens.Token{tokens.New(tokentype.LESS_LESS, "<<", 0, 0)}},
		{name: "great", text: ">", expected: []tokens.Token{tokens.New(tokentype.GREAT, ">", 0, 0)}},
		{name: "great equal", text: ">=", expected: []tokens.Token{tokens.New(tokentype.GREAT_EQUAL, ">=", 0, 0)}},
		{name: "great great", text: ">>", expected: []tokens.Token{tokens.New(tokentype.GREAT_GREAT, ">>", 0, 0)}},
		{name: "equal", text: "=", expected: []tokens.Token{tokens.New(tokentype.EQUAL, "=", 0, 0)}},
		{name: "equal equal", text: "==", expected: []tokens.Token{tokens.New(tokentype.EQUAL_EQUAL, "==", 0, 0)}},
		{name: "bang", text: "!", expected: []tokens.Token{tokens.New(tokentype.BANG, "!", 0, 0)}},
		{name: "bang equal", text: "!=", expected: []tokens.Token{tokens.New(tokentype.BANG_EQUAL, "!=", 0, 0)}},
		{name: "stroke", text: "|", expected: []tokens.Token{tokens.New(tokentype.STROKE, "|", 0, 0)}},
		{name: "stroke stroke", text: "||", expected: []tokens.Token{tokens.New(tokentype.STROKE_STROKE, "||", 0, 0)}},
		{name: "stroke great", text: "|>", expected: []tokens.Token{tokens.New(tokentype.STROKE_GREAT, "|>", 0, 0)}},
		{name: "ampersand", text: "&", expected: []tokens.Token{tokens.New(tokentype.AMPERSAND, "&", 0, 0)}},
		{name: "double ampersand", text: "&&", expected: []tokens.Token{tokens.New(tokentype.AMPERSAND_AMPERSAND, "&&", 0, 0)}},
		{name: "caret", text: "^", expected: []tokens.Token{tokens.New(tokentype.CARET, "^", 0, 0)}},
		{name: "tilde", text: "~", expected: []tokens.Token{tokens.New(tokentype.TILDE, "~", 0, 0)}},
		{name: "asterisk", text: "*", expected: []tokens.Token{tokens.New(tokentype.ASTERISK, "*", 0, 0)}},
		{name: "double asterisk", text: "**", expected: []tokens.Token{tokens.New(tokentype.ASTERISK_ASTERISK, "**", 0, 0)}},
		{name: "slash", text: "/", expected: []tokens.Token{tokens.New(tokentype.SLASH, "/", 0, 0)}},
		{name: "plus", text: "+", expected: []tokens.Token{tokens.New(tokentype.PLUS, "+", 0, 0)}},
		{name: "minus", text: "-", expected: []tokens.Token{tokens.New(tokentype.MINUS, "-", 0, 0)}},
		{name: "minus great", text: "->", expected: []tokens.Token{tokens.New(tokentype.MINUS_GREAT, "->", 0, 0)}},
		{name: "question", text: "?", expected: []tokens.Token{tokens.New(tokentype.QUESTION, "?", 0, 0)}},
		{name: "underscore", text: "_", expected: []tokens.Token{tokens.New(tokentype.UNDERSCORE, "_", 0, 0)}},
		{name: "number 1", text: "123", expected: []tokens.Token{tokens.New(tokentype.NUMBER, "123", 0, 0)}},
		{name: "number 2", text: "123.5", expected: []tokens.Token{tokens.New(tokentype.NUMBER, "123.5", 0, 0)}},
		{name: "number 3", text: "123.", expected: []tokens.Token{}, willError: true},
		{name: "identifier", text: "abc", expected: []tokens.Token{tokens.New(tokentype.IDENTIFIER, "abc", 0, 0)}},
		{name: "string double quotes", text: "\"abc\"", expected: []tokens.Token{tokens.New(tokentype.STRING, "\"abc\"", 0, 0)}},
		{name: "string single quotes", text: "'abc'", expected: []tokens.Token{tokens.New(tokentype.STRING, "'abc'", 0, 0)}},
		{name: "unterminated double string", text: "\"abc", expected: []tokens.Token{}, willError: true},
		{name: "unterminated single string", text: "'abc", expected: []tokens.Token{}, willError: true},
		{name: "if", text: "if", expected: []tokens.Token{tokens.New(tokentype.IF, "if", 0, 0)}},
		{name: "else", text: "else", expected: []tokens.Token{tokens.New(tokentype.ELSE, "else", 0, 0)}},
		{name: "for", text: "for", expected: []tokens.Token{tokens.New(tokentype.FOR, "for", 0, 0)}},
		{name: "let", text: "let", expected: []tokens.Token{tokens.New(tokentype.LET, "let", 0, 0)}},
		{name: "true", text: "true", expected: []tokens.Token{tokens.New(tokentype.TRUE, "true", 0, 0)}},
		{name: "false", text: "false", expected: []tokens.Token{tokens.New(tokentype.FALSE, "false", 0, 0)}},
		{name: "fn", text: "fn", expected: []tokens.Token{tokens.New(tokentype.FN, "fn", 0, 0)}},
		{name: "return", text: "return", expected: []tokens.Token{tokens.New(tokentype.RETURN, "return", 0, 0)}},
		{name: "bottom", text: "bottom", expected: []tokens.Token{tokens.New(tokentype.BOTTOM, "bottom", 0, 0)}},
		{name: "space", text: " ", expected: []tokens.Token{}},
		{name: "newline", text: "\n", expected: []tokens.Token{}},
		{name: "unrecognized symbol", text: "$", expected: []tokens.Token{}, willError: true},
		{name: "mut", text: "mut", expected: []tokens.Token{tokens.New(tokentype.MUT, "mut", 0, 0)}},
	}

	for _, e := range table {
		sc := scanner.New()
		ts, err := sc.Read(e.text)

		// Remove EOF token as it contributes nothing to the test
		ts = lessEOF(ts)

		if err != nil && !e.willError {
			t.Errorf("%s: Expected to receive an error when lexing %q", e.name, e.text)
		}

		if !same(ts, e.expected) {
			t.Errorf("%s: Tokens were not generated properly: got %#v instead of %#v", e.name, ts, e.expected)
		}
	}

	/*
		Produce multi-char tokens over single-char
	*/
	table = []struct {
		name      string
		text      string
		expected  []tokens.Token
		willError bool
	}{
		{
			name:     "triple less",
			text:     "<<<",
			expected: []tokens.Token{tokens.New(tokentype.LESS_LESS, "<<", 0, 0), tokens.New(tokentype.LESS, "<", 0, 0)},
		},
		{
			name:     "double stroke stroke great",
			text:     "|||>",
			expected: []tokens.Token{tokens.New(tokentype.STROKE_STROKE, "||", 0, 0), tokens.New(tokentype.STROKE_GREAT, "|>", 0, 0)},
		},
	}

	for _, e := range table {
		sc := scanner.New()
		ts, _ := sc.Read(e.text)
		ts = lessEOF(ts)

		if !same(ts, e.expected) {
			t.Errorf("%s: Did not parse multi-char tokens properly. Got %#v instead of %#v", e.name, ts, e.expected)
		}
	}

	/*
		Track positional information
	*/
	table = []struct {
		name      string
		text      string
		expected  []tokens.Token
		willError bool
	}{
		{
			name:     "column increments",
			text:     "+ -  <",
			expected: []tokens.Token{tokens.New(tokentype.PLUS, "+", 0, 0), tokens.New(tokentype.MINUS, "-", 0, 2), tokens.New(tokentype.LESS, "<", 0, 5)},
		},
		{
			name:     "row increments",
			text:     "+\n+",
			expected: []tokens.Token{tokens.New(tokentype.PLUS, "+", 0, 0), tokens.New(tokentype.PLUS, "+", 1, 0)},
		},
	}

	for _, e := range table {
		sc := scanner.New()
		ts, _ := sc.Read(e.text)
		ts = lessEOF(ts)

		if !samePos(ts, e.expected) {
			t.Errorf("%s: Positional information was wrong from lexer. Got %#v but expected %#v", e.name, ts, e.expected)
		}
	}
}
