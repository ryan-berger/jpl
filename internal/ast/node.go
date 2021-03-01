package ast

type SExpr interface {
	SExpr() string
}

type Locationer interface {
	Loc() (line, char int)
}

type Node interface {
	SExpr
	Locationer
	String() string
}
