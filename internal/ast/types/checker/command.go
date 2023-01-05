package checker

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/types"
	"github.com/ryan-berger/jpl/internal/symbol"
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

		srcType, err := expressionType(cmd.Src, table)
		if err != nil {
			return err
		}

		if !srcType.Equal(types.Str) {
			return NewError(cmd.Src, "type error: source in read command must be type str received: %s", srcType)
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

		dstType, err := expressionType(cmd.Dest, table)
		if err != nil {
			return err
		}

		if !dstType.Equal(types.Str) {
			return NewError(cmd.Dest, "type error: destination in write command must be type str received: %s", dstType)
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
