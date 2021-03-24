package expander

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

type nexter func() string

func defaultNexter() nexter {
	N := 0
	return func() string {
		N++
		return fmt.Sprintf("t.%d", N)
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
		name := next()
		l := dsl.Let(
			dsl.LIdent(name), dsl.Int(0))
		ret := dsl.Return(dsl.Ident(name))

		expanded = append(expanded, l, ret)
	}

	return expanded
}
