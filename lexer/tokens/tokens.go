package tokens

type Token struct {
	Type     tokentype
	Lexeme   string
	Position struct {
		Row int
		Col int
	}
}

func New(t tokentype, l string, r int, c int) Token {
	return Token{
		Type:   t,
		Lexeme: l,
		Position: struct {
			Row int
			Col int
		}{Row: r, Col: c},
	}
}
