package typed

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)


func bindArg(argument ast.Argument, typ types.Type, table *symbol.Table) error {
	switch arg := argument.(type) {
	case *ast.VariableArgument:
		if _, ok := table.Get(arg.Variable); ok {
			return NewError(arg,
				"cannot bindArg variable %s, variable is already bound", arg.Variable)
		}
		table.Set(arg.Variable, &symbol.Identifier{Type: typ})
		arg.Type = typ
		return nil
	case *ast.VariableArr:
		arrTyp, ok := typ.(*types.Array)
		if !ok {
			return NewError(arg, "array bindArg to non-array type for binding %s", arg.Variable)
		}
		if len(arg.Variables) != arrTyp.Rank {
			return NewError(arg, "dimension incorrect for binding %s", arg)
		}

		table.Set(arg.Variable, &symbol.Identifier{Type: arrTyp})
		for _, v := range arg.Variables {
			table.Set(v, &symbol.Identifier{Type: types.Integer})
		}
	}
	return nil
}

func bind(binding ast.Binding, table *symbol.Table) (types.Type, error) {
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
		if err := bindArg(b.Argument, b.Type, table); err != nil {
			return nil, err
		}
		return b.Type, nil
	default:
		panic("")
	}
	return nil, nil
}

func functionBinding(fun *ast.Function, table *symbol.Table) (*symbol.Function, error) {
	function := &symbol.Function{
		Args:   make([]types.Type, len(fun.Bindings)),
		Return: fun.ReturnType,
	}
	cpy := table.Copy()
	if _, ok := cpy.Get(fun.Var); ok {
		return nil, fmt.Errorf("error, function name %s already declared", fun.Var)
	}

	for i, b := range fun.Bindings {
		typ, err := bind(b, cpy)
		if err != nil {
			return nil, err
		}
		function.Args[i] = typ
	}

	hasSeenReturn := false
	for _, stmt := range fun.Statements {
		err := statementType(stmt, function.Return, table)
		if err != nil {
			return nil, err
		}

		_, isRet := stmt.(*ast.ReturnStatement)
		hasSeenReturn = hasSeenReturn || isRet
	}
	// return has not been found
	if !hasSeenReturn && !function.Return.Equal(&types.Tuple{}) {
		return nil, NewError(fun, "return of type expected %s, received none", function.Return)
	}

	return function, nil
}

func bindLVal(value ast.LValue, typ types.Type, table *symbol.Table) error {
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
	case ast.Argument:
		return bindArg(lval, typ, table)
	}
	return nil
}

func statementType(statement ast.Statement, retType types.Type, table *symbol.Table) error {
	switch stmt := statement.(type) {
	case *ast.Function:
		fn, err := functionBinding(stmt, table)
		if err != nil {
			return err
		}
		table.Set(stmt.Var, fn)
	case *ast.ReturnStatement:
		typ, err := expressionType(stmt.Expr, table)
		if err != nil {
			return err
		}

		if !typ.Equal(retType) {
			return NewError(stmt.Expr, "return type of %s expected, received %s", retType, typ)
		}
		return nil
	case *ast.LetStatement:
		rType, err := expressionType(stmt.Expr, table)
		if err != nil {
			return err
		}
		if err = bindLVal(stmt.LValue, rType, table); err != nil {
			return err
		}
	case *ast.AssertStatement:
		exprType, err := expressionType(stmt.Expr, table)
		if err != nil {
			return err
		}

		if !exprType.Equal(types.Boolean) {
			return fmt.Errorf("assert statement expression must be Boolean")
		}
	default:
		panic("statement not implemented")
	}
	return nil
}
