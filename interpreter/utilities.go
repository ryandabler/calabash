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
	tokentype.LESS:              nil,
	tokentype.LESS_EQUAL:        nil,
	tokentype.GREAT:             nil,
	tokentype.GREAT_EQUAL:       nil,
}

func isNumericOp(op tokentype.Tokentype) bool {
	_, ok := numericOps[op]
	return ok
}

func isBooleanOp(op tokentype.Tokentype) bool {
	return op == tokentype.AMPERSAND_AMPERSAND || op == tokentype.STROKE_STROKE
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

func areBools(bs ...interface{}) bool {
	for _, n := range bs {
		_, ok := n.(value.VBoolean)

		if !ok {
			return false
		}
	}

	return true
}
