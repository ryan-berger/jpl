package typed

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)

func checkRead(argument ast.Argument, table symbol.Table) error {
	switch arg := argument.(type) {
	case *ast.VariableArgument:
		if _, ok := table[arg.Variable]; ok {
			return NewError(arg, "variable %s already bound", arg.Variable)
		}
		table[arg.Variable] = &symbol.Identifier{Type: types.Integer}
	case *ast.VariableArr:
		if len(arg.Variables) != 2 {
			return NewError(arg, "read variable binding must be of rank two")
		}
		if _, ok := table[arg.Variable]; ok {
			return NewError(arg, "variable %s already bound", arg.Variable)
		}
		table[arg.Variable] = &symbol.Identifier{Type: types.Integer}

		for _, v := range arg.Variables {
			if _, ok := table[v]; ok {
				return NewError(arg, "variable %s already bound", arg.Variable)
			}
			table[v] = &symbol.Identifier{Type: types.Integer}
		}
	}

	return nil
}

func checkCommand(command ast.Command, table symbol.Table) (bool, error) {
	switch cmd := command.(type) {
	case ast.Statement:
		return statementType(cmd, table)
	case *ast.Read:
		if cmd.Type != "image" {
			return false, NewError(cmd, "oops, read type not supported yet")
		}
		if err := checkRead(cmd.Argument, table); err != nil {
			return false, err
		}
		return false, nil
	case *ast.Write:
		typ, err := expressionType(cmd.Expr, table)
		if err != nil {
			return false, err
		}

		if !typ.Equal(types.Pict) {
			return false, NewError(cmd.Expr, "type error: write command expected %s received %s", types.Pict, typ)
		}
		return false, nil
	case *ast.Show:
		_, err := expressionType(cmd.Expr, table)
		return false, err
	case *ast.Print:
		return false, nil
	case *ast.Time:
		return checkCommand(cmd.Command, table)
	default:
		panic("not implemented for command type")
	}
	return false, nil
}
