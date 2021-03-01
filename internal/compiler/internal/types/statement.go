package types

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

func toType(typ ast.Type) Type {
	switch t := typ.(type) {
	case ast.BasicType:
		switch t {
		case ast.Int:
			return Integer
		case ast.Float:
			return Float
		case ast.Boolean:
			return Boolean
		}
	case *ast.ArrType:
		return &Array{Inner: toType(t.Type), Rank: t.Rank}
	case *ast.TupleType:
		types := make([]Type, len(t.Types))
		for i, typ := range t.Types {
			types[i] = toType(typ)
		}
		return &Tuple{
			Types: types,
		}
	default:
		panic("invalid type")
	}
	return nil
}

func bindArg(argument ast.Argument, typ Type, table SymbolTable) error {
	switch arg := argument.(type) {
	case *ast.VariableArgument:
		if _, ok := table[arg.Variable]; ok {
			return fmt.Errorf("cannot bindArg variable %s, variable is already bound",
				arg.Variable)
		}
		table[arg.Variable] = &Identifier{Type: typ}
		return nil
	case *ast.VariableArr:
		arrTyp, ok := typ.(*Array)
		if !ok {
			return fmt.Errorf("array bindArg to non-array type for binding %s", arg.Variable)
		}
		if len(arg.Variables) != arrTyp.Rank {
			return fmt.Errorf("dimension incorrect for binding %s", arg)
		}

		table[arg.Variable] = &Identifier{Type: arrTyp}
		for _, v := range arg.Variables {
			table[v] = &Identifier{Type: Integer}
		}
	}
	return nil
}

func bind(binding ast.Binding, table SymbolTable) (Type, error) {
	switch b := binding.(type) {
	case *ast.TupleBinding:
		tup := &Tuple{
			Types: make([]Type, len(b.Bindings)),
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

func functionBinding(fun *ast.Function, table SymbolTable) (*Function, error) {
	function := &Function{
		Args:   make([]Type, len(fun.Bindings)),
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
		ret, err := StatementType(stmt, table)
		if err != nil {
			return nil, err
		}

		if ret {
			retStmt := stmt.(*ast.ReturnStatement)
			retTyp, err := ExpressionType(retStmt.Expr, table)
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
	if !function.Return.Equal(&Tuple{}) {
		return nil, fmt.Errorf("found no return, expected return of a type")
	}

	return function, nil
}

func bindLVal(value ast.LValue, typ Type, table SymbolTable) error {
	switch lval := value.(type) {
	case *ast.LTuple:
		tup, ok := typ.(*Tuple)
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
		arr, ok := typ.(*Array)
		if !ok {
			return fmt.Errorf("type must be array")
		}
		if arr.Rank != len(lval.Variables) {
			return fmt.Errorf("rank incorrect for binding: %s", lval.Variable)
		}
		table[lval.Variable] = &Identifier{Type: arr}
		for _, v := range lval.Variables {
			if _, ok := table[v]; ok {
				return fmt.Errorf("symbol already bound %s", lval.Variable)
			}
			table[v] = &Identifier{Type: Integer}
		}
	case *ast.VariableArgument:
		if _, ok := table[lval.Variable]; ok {
			return fmt.Errorf("symbol already bound %s", lval.Variable)
		}
		table[lval.Variable] = &Identifier{Type: typ}
	}
	return nil
}

func StatementType(statement ast.Statement, table SymbolTable) (bool, error) {
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
		rType, err := ExpressionType(stmt.Expr, table)
		if err != nil {
			return false, nil
		}
		if err = bindLVal(stmt.LValue, rType, table); err != nil {
			return false, err
		}
	case *ast.AssertStatement:
		exprType, err := ExpressionType(stmt.Expr, table)
		if err != nil {
			return false, err
		}

		if !exprType.Equal(Boolean) {
			return false, fmt.Errorf("assert statement expression must be Boolean")
		}
	default:
		panic("statement not implemented")
	}
	return false, nil
}
