package types

import (
	"github.com/ryan-berger/jpl/internal/ast"
)

func checkRead(argument ast.Argument, table SymbolTable) error {
	switch argument.(type) {
	case *ast.VariableArgument:

	case *ast.VariableArr:

	}


	return nil
}

func CheckCommand(command ast.Command, table SymbolTable) (bool, error) {
	switch cmd := command.(type) {
	case ast.Statement:
		return StatementType(cmd, table)
	case *ast.Read:
		if cmd.Type != "image" {
			panic("oops, not supported yet")
		}

	case *ast.Write:

	case *ast.Show:
		_, err := ExpressionType(cmd.Expr, table)
		return false, err
	case *ast.Print:
		return false, nil
	case *ast.Time:
		return CheckCommand(cmd.Command, table)
	default:
		panic("not implemented for command type")
	}
	return false, nil
}
