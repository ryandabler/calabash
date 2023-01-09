package tokens

import "calabash/internal/tokentype"

type Token struct {
	Type     tokentype.Tokentype
	Lexeme   string
	Position struct {
		Row int
		Col int
	}
}

func New(t tokentype.Tokentype, l string, r int, c int) Token {
	return Token{
		Type:   t,
		Lexeme: l,
		Position: struct {
			Row int
			Col int
		}{Row: r, Col: c},
	}
}
