package checker

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
)

func Check(program ast.Program) (ast.Program, error) {
	table := symbol.NewSymbolTable()

	for _, cmd := range program {
		err := checkCommand(cmd, table)
		if err != nil {
			return nil, err
		}
	}

	return program, nil
}
