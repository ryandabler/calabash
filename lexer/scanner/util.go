package scanner

import "calabash/lexer/tokens"

func isDigit(r rune) bool {
	if r == -1 {
		return false
	}

	return '0' <= r && r <= '9'
}

func isAlpha(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}

var keywords map[string]tokens.Token = map[string]tokens.Token{
	"if":     tokens.New(tokens.IF, "", 0, 0),
	"else":   tokens.New(tokens.ELSE, "", 0, 0),
	"for":    tokens.New(tokens.FOR, "", 0, 0),
	"let":    tokens.New(tokens.LET, "", 0, 0),
	"true":   tokens.New(tokens.TRUE, "", 0, 0),
	"false":  tokens.New(tokens.FALSE, "", 0, 0),
	"fn":     tokens.New(tokens.FN, "", 0, 0),
	"return": tokens.New(tokens.RETURN, "", 0, 0),
	"bottom": tokens.New(tokens.BOTTOM, "", 0, 0),
}
