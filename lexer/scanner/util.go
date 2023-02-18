package scanner

import (
	"calabash/internal/tokentype"
	"calabash/lexer/tokens"
)

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
	"if":     tokens.New(tokentype.IF, "", 0, 0),
	"else":   tokens.New(tokentype.ELSE, "", 0, 0),
	"for":    tokens.New(tokentype.FOR, "", 0, 0),
	"let":    tokens.New(tokentype.LET, "", 0, 0),
	"true":   tokens.New(tokentype.TRUE, "", 0, 0),
	"false":  tokens.New(tokentype.FALSE, "", 0, 0),
	"fn":     tokens.New(tokentype.FN, "", 0, 0),
	"return": tokens.New(tokentype.RETURN, "", 0, 0),
	"bottom": tokens.New(tokentype.BOTTOM, "", 0, 0),
	"mut":    tokens.New(tokentype.MUT, "", 0, 0),
}
