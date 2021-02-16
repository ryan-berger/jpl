package ast

type SExpr interface {
	SExpr() string
}

type Node interface {
	String() string
	SExpr
	Loc() (line, char int)
}
