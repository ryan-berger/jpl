package typed

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)

func toType(typ ast.Type) types.Type {
	switch t := typ.(type) {
	case ast.BasicType:
		switch t {
		case ast.Int:
			return types.Integer
		case ast.Float:
			return types.Float
		case ast.Boolean:
			return types.Boolean
		}
	case *ast.ArrType:
		return &types.Array{Inner: toType(t.Type), Rank: t.Rank}
	case *ast.TupleType:
		typs := make([]types.Type, len(t.Types))
		for i, typ := range t.Types {
			typs[i] = toType(typ)
		}
		return &types.Tuple{
			Types: typs,
		}
	default:
		panic("invalid type")
	}
	return nil
}

func bindArg(argument ast.Argument, typ types.Type, table symbol.Table) error {
	switch arg := argument.(type) {
	case *ast.VariableArgument:
		if _, ok := table[arg.Variable]; ok {
			return fmt.Errorf("cannot bindArg variable %s, variable is already bound",
				arg.Variable)
		}
		table[arg.Variable] = &symbol.Identifier{Type: typ}
		return nil
	case *ast.VariableArr:
		arrTyp, ok := typ.(*types.Array)
		if !ok {
			return fmt.Errorf("array bindArg to non-array type for binding %s", arg.Variable)
		}
		if len(arg.Variables) != arrTyp.Rank {
			return fmt.Errorf("dimension incorrect for binding %s", arg)
		}

		table[arg.Variable] = &symbol.Identifier{Type: arrTyp}
		for _, v := range arg.Variables {
			table[v] = &symbol.Identifier{Type: types.Integer}
		}
	}
	return nil
}

func bind(binding ast.Binding, table symbol.Table) (types.Type, error) {
	switch b := binding.(type) {
	case *ast.TupleBinding:
		tup := &types.Tuple{
			Types: make([]types.Type, len(b.Bindings)),
		}
		for i, binding := range b.Bindings {
			typ, err := bind(binding, table)
			if err != nil {
				return nil, err
			}
			tup.Types[i] = typ
		}
		return tup, nil
	case *ast.TypeBind:
		typ := toType(b.Type)
		if err := bindArg(b.Argument, typ, table); err != nil {
			return nil, err
		}
		return typ, nil
	default:
		panic("")
	}
	return nil, nil
}

func functionBinding(fun *ast.Function, table symbol.Table) (*symbol.Function, error) {
	function := &symbol.Function{
		Args:   make([]types.Type, len(fun.Bindings)),
		Return: toType(fun.ReturnType),
	}
	cpy := table.Copy()
	if _, ok := cpy[fun.Var]; ok {
		return nil, fmt.Errorf("error, function name %s already declared", fun.Var)
	}

	for i, b := range fun.Bindings {
		typ, err := bind(b, cpy)
		if err != nil {
			return nil, err
		}
		function.Args[i] = typ
	}

	for _, stmt := range fun.Statements {
		ret, err := statementType(stmt, table)
		if err != nil {
			return nil, err
		}

		if ret {
			retStmt := stmt.(*ast.ReturnStatement)
			retTyp, err := expressionType(retStmt.Expr, table)
			if err != nil {
				return nil, err
			}

			if !retTyp.Equal(function.Return) {
				return nil, fmt.Errorf("function return expects different type")
			}
			return function, nil
		}
	}
	// return has not been found
	if !function.Return.Equal(&types.Tuple{}) {
		return nil, NewError(fun, "found no return, expected return of a type")
	}

	return function, nil
}

func bindLVal(value ast.LValue, typ types.Type, table symbol.Table) error {
	switch lval := value.(type) {
	case *ast.LTuple:
		tup, ok := typ.(*types.Tuple)
		if !ok {
			return fmt.Errorf("expected tuple binding")
		}

		if len(tup.Types) != len(lval.Args) {
			return fmt.Errorf("tuples are different shapes")
		}

		for i, v := range lval.Args {
			if err := bindLVal(v, tup.Types[i], table); err != nil {
				return err
			}
		}
	case *ast.VariableArr:
		arr, ok := typ.(*types.Array)
		if !ok {
			return fmt.Errorf("type must be array")
		}
		if arr.Rank != len(lval.Variables) {
			return fmt.Errorf("rank incorrect for binding: %s", lval.Variable)
		}
		table[lval.Variable] = &symbol.Identifier{Type: arr}
		for _, v := range lval.Variables {
			if _, ok := table[v]; ok {
				return fmt.Errorf("symbol already bound %s", lval.Variable)
			}
			table[v] = &symbol.Identifier{Type: types.Integer}
		}
	case *ast.VariableArgument:
		if _, ok := table[lval.Variable]; ok {
			return fmt.Errorf("symbol already bound %s", lval.Variable)
		}
		table[lval.Variable] = &symbol.Identifier{Type: typ}
	}
	return nil
}

func statementType(statement ast.Statement, table symbol.Table) (bool, error) {
	switch stmt := statement.(type) {
	case *ast.Function:
		fn, err := functionBinding(stmt, table)
		if err != nil {
			return false, err
		}
		table[stmt.Var] = fn
	case *ast.ReturnStatement:
		return true, nil
	case *ast.LetStatement:
		rType, err := expressionType(stmt.Expr, table)
		if err != nil {
			return false, nil
		}
		if err = bindLVal(stmt.LValue, rType, table); err != nil {
			return false, err
		}
	case *ast.AssertStatement:
		exprType, err := expressionType(stmt.Expr, table)
		if err != nil {
			return false, err
		}

		if !exprType.Equal(types.Boolean) {
			return false, fmt.Errorf("assert statement expression must be Boolean")
		}
	default:
		panic("statement not implemented")
	}
	return false, nil
}
