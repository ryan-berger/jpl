package typed

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)

func checkIf(ifExpr *ast.IfExpression, table *symbol.Table) (types.Type, error) {
	condType, err := expressionType(ifExpr.Condition, table)
	if err != nil {
		return nil, err
	}
	if !condType.Equal(types.Boolean) {
		return nil, NewError(ifExpr.Condition, "expected boolean, received %s", condType)
	}

	consType, err := expressionType(ifExpr.Consequence, table)
	if err != nil {
		return nil, err
	}

	otherType, err := expressionType(ifExpr.Otherwise, table)
	if err != nil {
		return nil, err
	}

	if !consType.Equal(otherType) {
		return nil, NewError(ifExpr, "branches return different types: %s, %s", consType, otherType)
	}

	ifExpr.Type = consType
	return otherType, nil
}

func checkIdentifierExpr(expr *ast.IdentifierExpression, table *symbol.Table) (types.Type, error) {
	symb, ok := table.Get(expr.Identifier)
	if !ok {
		return nil, NewError(expr, "unknown symbol %s", expr.Identifier)
	}

	ident, ok := symb.(*symbol.Identifier)
	if !ok {
		return nil, fmt.Errorf("found function name, expected identifier %s", expr.Identifier)
	}

	expr.Type = ident.Type
	return ident.Type, nil
}

func checkDimExpr(expr *ast.CallExpression, table *symbol.Table) (types.Type, error) {
	if len(expr.Arguments) != 2 {
		return nil, NewError(
			expr,
			"function dim expects 2 arguments, received %d", len(expr.Arguments))
	}

	arrTyp, err := expressionType(expr.Arguments[0], table)
	if err != nil {
		return nil, err
	}

	if _, ok := arrTyp.(*types.Array); !ok {
		return nil, NewError(expr.Arguments[0], "type error: expected array type, received %s", arrTyp.String())
	}

	dimTyp, err := expressionType(expr.Arguments[1], table)
	if err != nil {
		return nil, err
	}

	if !dimTyp.Equal(types.Integer) {
		return nil, NewError(expr.Arguments[1], "type error: expected int received %s", dimTyp.String())
	}

	return types.Integer, nil
}

func checkCallExpr(expr *ast.CallExpression, table *symbol.Table) (types.Type, error) {
	// special case for dim
	if expr.Identifier == "dim" {
		return checkDimExpr(expr, table)
	}

	symb, ok := table.Get(expr.Identifier)
	if !ok {
		return nil, NewError(expr, "unknown symbol %s", expr.Identifier)
	}

	call, ok := symb.(*symbol.Function)
	if !ok {
		return nil, NewError(expr, "found identifier, expected function name %s", expr.Identifier)
	}

	if len(call.Args) != len(expr.Arguments) {
		return nil, NewError(expr, "function %s expects %d arguments, received %d", expr.Identifier, len(call.Args), len(expr.Arguments))
	}

	for i, t := range call.Args {
		exprType, err := expressionType(expr.Arguments[i], table)
		if err != nil {
			return nil, err
		}
		if !t.Equal(exprType) {
			return nil, NewError(expr.Arguments[i], "type error: expected %s received %s", t, exprType)
		}
	}

	expr.Type = call.Return
	return call.Return, nil
}

func isNumeric(typ types.Type) bool {
	return typ.Equal(types.Float) || typ.Equal(types.Integer)
}

func checkInfixExpr(expression *ast.InfixExpression, table *symbol.Table) (types.Type, error) {
	leftType, leftErr := expressionType(expression.Left, table)
	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := expressionType(expression.Right, table)
	if rightErr != nil {
		return nil, rightErr
	}

	// if the expression is an and/or, check for Boolean types
	if expression.Op == "&&" || expression.Op == "||" {
		if !leftType.Equal(types.Boolean) {
			return nil, NewError(expression.Left,
				"type error: left operand of %s expression is of type %s expected bool", expression.Op, leftType)
		}
		if !rightType.Equal(types.Boolean) {
			return nil, NewError(expression.Right,
				"type error: right operand of %s expression is of type %s expected bool", expression.Op, rightType)

		}

		expression.Type = types.Boolean
		return types.Boolean, nil
	}

	// if we don't have logical operators, we need to make sure that both
	// expressions are the same and are not Booleans
	if !isNumeric(leftType) {
		return nil, NewError(expression.Left,
			"type error: left type of %s expression must be numeric, received %s", expression.Op, leftType)
	}
	if !isNumeric(rightType) {
		return nil, NewError(expression.Left,
			"type error: right type of %s expression must be numeric, received %s", expression.Op, rightType)
	}

	if !leftType.Equal(rightType) {
		return nil, NewError(expression, "type error: both sides of numerical operation must be of the same type")
	}

	switch expression.Op {
	// comparison operators
	case "==", "!=", "<=", ">=", "<", ">":
		expression.Type = types.Boolean
		return types.Boolean, nil
	case "+", "-", "*", "/", "%":
		expression.Type = leftType
		return leftType, nil
	default:
		panic("invalid infix operation")
	}
	return nil, nil
}

func checkPrefixExpr(expr *ast.PrefixExpression, table *symbol.Table) (types.Type, error) {
	exprType, err := expressionType(expr.Expr, table)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case "!":
		if !exprType.Equal(types.Boolean) {
			return nil, NewError(expr,
				"type error, expected boolean on right hand side of '!', received: %s", exprType)
		}
		expr.Type = types.Boolean
		return types.Boolean, nil
	case "-":
		if !isNumeric(exprType) {
			return nil, NewError(expr,
				"type error, expected numeric type on right hand side of '-', received: %s", exprType)
		}
		expr.Type = exprType
		return exprType, nil
	default:
		panic("invalid prefix operation")
	}
}

func checkTupleRef(expr *ast.TupleRefExpression, table *symbol.Table) (types.Type, error) {
	tup, err := expressionType(expr.Tuple, table)
	if err != nil {
		return nil, err
	}

	tupType, ok := tup.(*types.Tuple)
	if !ok {
		return nil, NewError(expr.Tuple, "tuple index of non-tuple type %s", tup)
	}

	if expr.Index < 0 || int(expr.Index) > len(tupType.Types)-1 {
		return nil, NewError(expr, "tuple index out of bounds")
	}

	return tupType.Types[expr.Index], nil
}

func checkArrayRef(expr *ast.ArrayRefExpression, table *symbol.Table) (types.Type, error) {
	arr, err := expressionType(expr.Array, table)
	if err != nil {
		return nil, err
	}

	arrTyp, ok := arr.(*types.Array)
	if !ok {
		return nil, NewError(expr.Array, "array reference of non-array type %s", arr)
	}

	if arrTyp.Rank != len(expr.Indexes) {
		return nil, NewError(expr.Array, "array access of rank %d with %d indexes",
			arrTyp.Rank, len(expr.Indexes))
	}

	for _, idxExp := range expr.Indexes {
		idxTyp, err := expressionType(idxExp, table)
		if err != nil {
			return nil, err
		}
		if !idxTyp.Equal(types.Integer) {
			return nil, NewError(idxExp, "non-integer index of array type %s", idxTyp)
		}
	}

	return arrTyp.Inner, nil
}

func checkSumTransform(expr *ast.SumTransform, table *symbol.Table) (types.Type, error) {
	cpy := table.Copy()
	for _, binding := range expr.OpBindings {
		if _, ok := cpy.Get(binding.Variable); ok {
			return nil, NewError(expr, "illegal shadowing in sum expr, var: %s", binding.Variable)
		}
		bindType, err := expressionType(binding.Expr, cpy)
		if err != nil {
			return nil, err
		}

		if !bindType.Equal(types.Integer) {
			return nil, NewError(binding.Expr, "bindArg expr initializer for %s returns non-integer",
				binding.Variable)
		}

		cpy.Set(binding.Variable, &symbol.Identifier{Type: types.Integer})
	}
	exprType, err := expressionType(expr.Expr, cpy)
	if err != nil {
		return nil, err
	}

	if !isNumeric(exprType) {
		return nil, NewError(expr.Expr, "sum returns non-numeric expression")
	}

	return exprType, nil
}

func checkArrayTransform(expr *ast.ArrayTransform, table *symbol.Table) (types.Type, error) {
	// copy symbol table to use as a local copy
	cpy := table.Copy()

	// loop over bindings and
	for _, binding := range expr.OpBindings {
		// make sure no variable shadowing is going on
		if _, ok := cpy.Get(binding.Variable); ok {
			return nil, NewError(binding,
				"illegal shadowing in sum expr, var: %s", binding.Variable)
		}

		// typecheck the binding type <var> : <expression> <-- this here
		bindType, err := expressionType(binding.Expr, cpy)
		if err != nil {
			return nil, err
		}

		// the binding type must be an integer
		if !bindType.Equal(types.Integer) {
			return nil, NewError(binding, "bindArg expr initializer for %s returns non-integer",
				binding.Variable)
		}

		// set the variable up in the local symbol table as an integer
		cpy.Set(binding.Variable, &symbol.Identifier{Type: types.Integer})
	}

	// get the expression type of the right hand side of the expression
	exprType, err := expressionType(expr.Expr, cpy)
	if err != nil {
		return nil, err
	}


	rank := len(expr.OpBindings)
	var typ types.Type = &types.Array{
		Inner: exprType,
		Rank: rank,
	}

	if rank == 0 {
		typ = exprType
	}

	expr.Type = typ
	return typ, nil
}

func checkTuple(expr *ast.TupleExpression, table *symbol.Table) (types.Type, error) {
	tuple := &types.Tuple{
		Types: make([]types.Type, len(expr.Expressions)),
	}
	for i, expr := range expr.Expressions {
		typ, err := expressionType(expr, table)
		if err != nil {
			return nil, err
		}
		tuple.Types[i] = typ
	}

	expr.Type = tuple
	return tuple, nil
}

func checkArray(expr *ast.ArrayExpression, table *symbol.Table) (types.Type, error) {
	if len(expr.Expressions) == 0 {
		return &types.Array{Inner: types.Integer, Rank: 1}, nil
	}

	typ, err := expressionType(expr.Expressions[0], table)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(expr.Expressions); i++ {
		curTyp, err := expressionType(expr.Expressions[i], table)
		if err != nil {
			return nil, err
		}

		if !curTyp.Equal(typ) {
			return nil, NewError(expr, "array literal has mixed types")
		}
	}

	return &types.Array{Inner: typ, Rank: 1}, nil
}

func expressionType(expression ast.Expression, table *symbol.Table) (types.Type, error) {
	switch expr := expression.(type) {
	case *ast.BooleanExpression:
		expr.Type = types.Boolean
		return types.Boolean, nil
	case *ast.IntExpression:
		expr.Type = types.Integer
		return types.Integer, nil
	case *ast.FloatExpression:
		expr.Type = types.Float
		return types.Float, nil
	case *ast.IdentifierExpression:
		return checkIdentifierExpr(expr, table)
	case *ast.CallExpression:
		return checkCallExpr(expr, table)
	case *ast.InfixExpression:
		return checkInfixExpr(expr, table)
	case *ast.PrefixExpression:
		return checkPrefixExpr(expr, table)
	case *ast.TupleRefExpression:
		return checkTupleRef(expr, table)
	case *ast.ArrayRefExpression:
		return checkArrayRef(expr, table)
	case *ast.IfExpression:
		return checkIf(expr, table)
	case *ast.SumTransform:
		return checkSumTransform(expr, table)
	case *ast.ArrayTransform:
		return checkArrayTransform(expr, table)
	case *ast.TupleExpression:
		return checkTuple(expr, table)
	case *ast.ArrayExpression:
		return checkArray(expr, table)
	default:
		panic(fmt.Sprintf("typechecking not implemented for type %T", expr))

	}
	return nil, nil
}
