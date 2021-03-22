package expander

import (
	"github.com/ryan-berger/jpl/internal/ast"
)

type nexter func() int

func defaultNexter() nexter {
	N := 0
	return func() int {
		N++
		return N
	}
}

func isReturn(cmd ast.Command) bool {
	_, ok := cmd.(*ast.ReturnStatement)
	return ok
}

func Expand(program ast.Program) ast.Program {
	var expanded ast.Program

	next := defaultNexter()

	for _, cmd := range program {
		expanded = append(expanded, expandCommand(cmd, next)...)

		if size := len(expanded); size != 0 {
			last := expanded[size-1]
			if _, ok := last.(*ast.ReturnStatement); ok { // we are done since we hit return
				return expanded
			}
		}
	}

	if size := len(expanded);
		size == 0 || !isReturn(expanded[size-1]) { // add a return at last since there is none
		l := let(next(), constInt(0))
		ret := returnStmt(refExpr(ident(l.LValue)))
		expanded = append(expanded, l, ret)
	}

	return expanded
}
