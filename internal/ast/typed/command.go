package typed

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)

func checkRead(argument ast.Argument, table *symbol.Table) error {
	return bindArg(argument, types.Pict, table)
}

func checkCommand(command ast.Command, table *symbol.Table) error {
	switch cmd := command.(type) {
	case ast.Statement:
		return statementType(cmd, types.Integer, table)
	case *ast.Function:
		return functionBinding(cmd, table)
	case *ast.Read:
		if cmd.Type != "image" {
			return NewError(cmd, "oops, read type not supported yet")
		}
		return checkRead(cmd.Argument, table)
	case *ast.Write:
		typ, err := expressionType(cmd.Expr, table)
		if err != nil {
			return err
		}

		if !typ.Equal(types.Pict) {
			return NewError(cmd.Expr, "type error: write command expected %s received %s", types.Pict, typ)
		}
		return nil
	case *ast.Show:
		_, err := expressionType(cmd.Expr, table)
		return err
	case *ast.Print:
		return nil
	case *ast.Time:
		return checkCommand(cmd.Command, table)
	default:
		panic("not implemented for command type")
	}
	return nil
}
