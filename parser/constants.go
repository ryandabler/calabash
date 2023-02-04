package parser

import (
	"calabash/internal/tokentype"
)

var comparisonTokens []tokentype.Tokentype = []tokentype.Tokentype{
	tokentype.GREAT,
	tokentype.GREAT_EQUAL,
	tokentype.LESS,
	tokentype.LESS_EQUAL,
}

var equalityTokens []tokentype.Tokentype = []tokentype.Tokentype{
	tokentype.EQUAL_EQUAL,
	tokentype.BANG_EQUAL,
}
