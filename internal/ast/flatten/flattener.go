package flatten

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

func Flatten(program ast.Program) ast.Program {
	var flattened ast.Program

	next := defaultNexter()

	for _, cmd := range program {
		flattened = append(flattened, flattenCommand(cmd, next)...)

		if size := len(flattened); size != 0 {
			last := flattened[size-1]
			if _, ok := last.(*ast.ReturnStatement); ok { // we are done since we hit return
				return flattened
			}
		}
	}

	if size := len(flattened);
		size == 0 || !isReturn(flattened[size-1]) { // add a return at last since there is none
		name := next()
		l := dsl.Let(
			dsl.LIdent(name), dsl.Int(0))
		ret := dsl.Return(dsl.Ident(name))

		flattened = append(flattened, l, ret)
	}

	return flattened
}
