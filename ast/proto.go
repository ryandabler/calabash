package ast

type ProtoMethod struct {
	K Expr // Method name
	M Expr // Method
	I bool // Auto-inherits
}
