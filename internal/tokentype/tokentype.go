package tokentype

type Tokentype int

const (
	LEFT_PAREN Tokentype = iota
	RIGHT_PAREN
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	SEMICOLON
	LESS
	LESS_EQUAL
	LESS_LESS
	GREAT
	GREAT_EQUAL
	GREAT_GREAT
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	STROKE
	STROKE_STROKE
	STROKE_GREAT
	AMPERSAND
	AMPERSAND_AMPERSAND
	CARET
	TILDE
	ASTERISK
	ASTERISK_ASTERISK
	SLASH
	PLUS
	MINUS
	MINUS_GREAT
	QUESTION
	UNDERSCORE
	NUMBER
	IDENTIFIER
	STRING
	IF
	ELSE
	FOR
	LET
	TRUE
	FALSE
	FN
	RETURN
	BOTTOM
	MUT
	ME
	PROTO
	WHILE
	CONTINUE
	BREAK
	DOT_DOT_DOT
	EOF
)
