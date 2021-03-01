package ast

type Location struct {
	Line, Pos int
}

func (l Location) Loc() (int, int) {
	return l.Line, l.Pos
}
