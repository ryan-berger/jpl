package types

import (
	"fmt"

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
		isReturn, err := StatementType(cmd, table)
		if err != nil {}
		fmt.Println(isReturn, err)
	case *ast.Read:
		if cmd.Type != "image" {
			return false, NewError(cmd, "oops, read type not supported yet")
		}
		
		return false, nil
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
