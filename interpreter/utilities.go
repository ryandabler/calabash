package interpreter

import (
	"calabash/internal/tokentype"
	"calabash/internal/value"
)

var numericOps map[tokentype.Tokentype]interface{} = map[tokentype.Tokentype]interface{}{
	tokentype.MINUS:             nil,
	tokentype.ASTERISK:          nil,
	tokentype.SLASH:             nil,
	tokentype.ASTERISK_ASTERISK: nil,
}

func isNumericOp(op tokentype.Tokentype) bool {
	_, ok := numericOps[op]
	return ok
}

func areNumbers(ns ...interface{}) bool {
	for _, n := range ns {
		_, ok := n.(value.VNumber)

		if !ok {
			return false
		}
	}

	return true
}
