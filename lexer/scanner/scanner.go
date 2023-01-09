package scanner

import (
	"calabash/errors"
	"calabash/internal/tokentype"
	"calabash/lexer/tokens"
	"fmt"
)

type scanner struct {
	rs  []rune
	cur int
	pos struct {
		col int
		row int
	}
}

func (s *scanner) Read(str string) ([]tokens.Token, error) {
	s.rs = []rune(str)
	s.cur = 0
	ts := []tokens.Token{}

	for !s.isEnd() {
		switch s.char() {
		case ' ':

		case '\n':
			s.pos.row++
			s.pos.col = -1

		case '(':
			ts = append(ts, tokens.New(tokentype.LEFT_PAREN, "(", s.pos.row, s.pos.col))

		case ')':
			ts = append(ts, tokens.New(tokentype.RIGHT_PAREN, ")", s.pos.row, s.pos.col))

		case '[':
			ts = append(ts, tokens.New(tokentype.LEFT_BRACKET, "[", s.pos.row, s.pos.col))

		case ']':
			ts = append(ts, tokens.New(tokentype.RIGHT_BRACKET, "]", s.pos.row, s.pos.col))

		case '{':
			ts = append(ts, tokens.New(tokentype.LEFT_BRACE, "{", s.pos.row, s.pos.col))

		case '}':
			ts = append(ts, tokens.New(tokentype.RIGHT_BRACE, "}", s.pos.row, s.pos.col))

		case ',':
			ts = append(ts, tokens.New(tokentype.COMMA, ",", s.pos.row, s.pos.col))

		case ';':
			ts = append(ts, tokens.New(tokentype.SEMICOLON, ";", s.pos.row, s.pos.col))

		case '?':
			ts = append(ts, tokens.New(tokentype.QUESTION, "?", s.pos.row, s.pos.col))

		case '_':
			ts = append(ts, tokens.New(tokentype.UNDERSCORE, "_", s.pos.row, s.pos.col))

		case '<':
			{
				next := s.peek()

				if next == '=' {
					ts = append(ts, tokens.New(tokentype.LESS_EQUAL, "<=", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else if next == '<' {
					ts = append(ts, tokens.New(tokentype.LESS_LESS, "<<", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.LESS, "<", s.pos.row, s.pos.col))
				}
			}

		case '>':
			{
				next := s.peek()

				if next == '=' {
					ts = append(ts, tokens.New(tokentype.GREAT_EQUAL, ">=", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else if next == '>' {
					ts = append(ts, tokens.New(tokentype.GREAT_GREAT, ">>", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.GREAT, ">", s.pos.row, s.pos.col))
				}
			}

		case '=':
			{
				next := s.peek()

				if next == '=' {
					ts = append(ts, tokens.New(tokentype.EQUAL_EQUAL, "==", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.EQUAL, "=", s.pos.row, s.pos.col))
				}
			}

		case '!':
			{
				next := s.peek()

				if next == '=' {
					ts = append(ts, tokens.New(tokentype.BANG_EQUAL, "!=", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.BANG, "!", s.pos.row, s.pos.col))
				}
			}

		case '|':
			{
				next := s.peek()

				if next == '|' {
					ts = append(ts, tokens.New(tokentype.STROKE_STROKE, "||", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else if next == '>' {
					ts = append(ts, tokens.New(tokentype.STROKE_GREAT, "|>", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.STROKE, "|", s.pos.row, s.pos.col))
				}
			}

		case '&':
			{
				next := s.peek()

				if next == '&' {
					ts = append(ts, tokens.New(tokentype.AMPERSAND_AMPERSAND, "&&", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.AMPERSAND, "&", s.pos.row, s.pos.col))
				}
			}

		case '^':
			ts = append(ts, tokens.New(tokentype.CARET, "^", s.pos.row, s.pos.col))

		case '~':
			ts = append(ts, tokens.New(tokentype.TILDE, "~", s.pos.row, s.pos.col))

		case '/':
			ts = append(ts, tokens.New(tokentype.SLASH, "/", s.pos.row, s.pos.col))

		case '*':
			{
				next := s.peek()

				if next == '*' {
					ts = append(ts, tokens.New(tokentype.ASTERISK_ASTERISK, "**", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.ASTERISK, "*", s.pos.row, s.pos.col))
				}
			}

		case '+':
			ts = append(ts, tokens.New(tokentype.PLUS, "+", s.pos.row, s.pos.col))

		case '-':
			{
				next := s.peek()

				if next == '>' {
					ts = append(ts, tokens.New(tokentype.MINUS_GREAT, "->", s.pos.row, s.pos.col))
					s.next() // Move ahead one token since we have a two-character token
				} else {
					ts = append(ts, tokens.New(tokentype.MINUS, "-", s.pos.row, s.pos.col))
				}
			}

		case '"':
			{
				cs := []rune{s.char()}
				col := s.pos.col
				s.next()

				for s.char() != '"' {
					if s.isEnd() {
						return []tokens.Token{}, errors.ScanError{Msg: "Unterminated string literal"}
					}

					cs = append(cs, s.char())
					s.next()
				}

				cs = append(cs, s.char())
				s.next()

				ts = append(ts, tokens.New(tokentype.STRING, string(cs), s.pos.row, col))
			}

		case '\'':
			{
				cs := []rune{s.char()}
				col := s.pos.col
				s.next()

				for s.char() != '\'' {
					if s.isEnd() {
						return []tokens.Token{}, errors.ScanError{Msg: "Unterminated string literal"}
					}

					cs = append(cs, s.char())
					s.next()
				}

				cs = append(cs, s.char())
				s.next()

				ts = append(ts, tokens.New(tokentype.STRING, string(cs), s.pos.row, col))
			}

		default:
			tl := len(ts)

			if isDigit(s.char()) {
				ds := []rune{s.char()}
				col := s.pos.col

				for s.peek() != '.' && isDigit(s.peek()) {
					s.next()
					ds = append(ds, s.char())
				}

				if s.peek() == '.' {
					s.next()
					ds = append(ds, s.char())

					if !isDigit(s.peek()) {
						return []tokens.Token{}, errors.ScanError{Msg: "Decimals must have digits after the decimal point"}
					}
				}

				for isDigit(s.peek()) {
					s.next()
					ds = append(ds, s.char())
				}

				ts = append(ts, tokens.New(tokentype.NUMBER, string(ds), s.pos.row, col))
			}

			if isAlpha(s.char()) {
				as := []rune{s.char()}
				col := s.pos.col

				for isAlpha(s.peek()) {
					s.next()
					as = append(as, s.char())
				}

				word := string(as)

				tk, ok := keywords[word]

				// If word is not a keyword it must be an identifier so create.
				// Else replace dummy data with actual data from scanner for
				// keyword token.
				if !ok {
					tk = tokens.New(tokentype.IDENTIFIER, word, s.pos.row, col)
				} else {
					tk.Lexeme = word
					tk.Position.Row = s.pos.row
					tk.Position.Col = col
				}

				ts = append(ts, tk)
			}

			if tl == len(ts) {
				return []tokens.Token{}, errors.ScanError{Msg: fmt.Sprintf("Unrecognizable symbol %q at (%d, %d)", s.char(), s.pos.row, s.pos.col)}
			}
		}

		s.next()
	}

	return ts, nil
}

func (s *scanner) isEnd() bool {
	return s.cur >= len(s.rs)
}

func (s *scanner) next() {
	// In addition to advancing the cursor we also need to advance our scanner's position
	// so that future tokens will have correct positional information
	s.pos.col++
	s.cur++
}

func (s *scanner) char() rune {
	if s.isEnd() {
		return -1
	}

	return s.rs[s.cur]
}

func (s *scanner) peek() rune {
	if s.cur+1 >= len(s.rs) {
		return -1
	}

	return s.rs[s.cur+1]
}

func New() *scanner {
	return &scanner{}
}
