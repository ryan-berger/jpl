package typed

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
)

func Check(program ast.Program) (ast.Program, error) {
	table := symbol.NewSymbolTable()

	i := 0
	for i < len(program) {
		isRet, err := checkCommand(program[i], table)
		if err != nil {
			return nil, err
		}
		if isRet {
			break
		}
		i++
	}

	return program[:i], nil
}
