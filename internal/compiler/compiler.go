package compiler

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/compiler/internal/types"
)

func typeCheckProgram(commands []ast.Command) error {
	symb := types.NewSymbolTable()
	for _, command := range commands {
		switch cmd := command.(type) {

		case ast.Statement:
			ret, err := types.StatementType(cmd, symb)
			if err != nil {
				return err
			}

			if ret {
				return nil
			}
		}
	}
	return nil
}

func Compile(ast []ast.Command) {}
