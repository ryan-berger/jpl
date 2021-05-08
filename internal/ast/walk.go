package ast


// TODO: is this even useful?
type WalkFn func(node Node)
func Walk(node Node, f WalkFn) {
	f(node)
	switch n := node.(type) {
	case Program:
		for _, cmd := range n {
			f(cmd)
		}
	case *LetStatement:

	}
}
