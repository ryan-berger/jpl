package expander

import "github.com/ryan-berger/jpl/internal/ast"

type nexter func() int

func defaultNexter() nexter {
	N := 0
	return func() int {
		N++
		return N
	}
}

func Expand(program ast.Program) ast.Program {
	var expanded ast.Program

	next := defaultNexter()

	for _, cmd := range program {
		expanded = append(expanded, expandCommand(cmd, next)...)
	}

	return expanded
}