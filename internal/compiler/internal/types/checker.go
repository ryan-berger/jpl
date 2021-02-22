package types

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

func checkIf(ifExpr *ast.IfExpression, table SymbolTable) (Type, error) {
	condType, err := ExpressionType(ifExpr.Condition, table)
	if err != nil {
		return nil, err
	}
	if !condType.Equal(Boolean) {
		return nil, fmt.Errorf("conditional expression is not boolean")
	}

	consType, err := ExpressionType(ifExpr.Consequence, table)
	if err != nil {
		return nil, err
	}

	otherType, err := ExpressionType(ifExpr.Otherwise, table)
	if err != nil {
		return nil, err
	}

	if !consType.Equal(otherType) {
		return nil, fmt.Errorf("branches do not return same type")
	}

	return otherType, nil
}

func checkIdentifierExpr(expr *ast.IdentifierExpression, table SymbolTable) (Type, error) {
	symb, ok := table[expr.Identifier]
	if !ok {
		return nil, fmt.Errorf("unknown symbol %s", expr.Identifier)
	}

	ident, ok := symb.(*Identifier)
	if !ok {
		return nil, fmt.Errorf("found function name, expected identifier %s", expr.Identifier)
	}

	return ident.Type, nil
}

func checkCallExpr(expr *ast.CallExpression, table SymbolTable) (Type, error) {
	symb, ok := table[expr.Identifier]
	if !ok {
		return nil, fmt.Errorf("unknown symbol %s", expr.Identifier)
	}

	call, ok := symb.(*Function)
	if !ok {
		return nil, fmt.Errorf("found identifier, expected function name %s", expr.Identifier)
	}

	if len(call.Args) != len(expr.Arguments) {
		return nil, fmt.Errorf("function %s expects %d arguments, received %d", expr.Identifier, len(call.Args), len(expr.Arguments))
	}

	for i, t := range call.Args {
		exprType, err := ExpressionType(expr.Arguments[i], table)
		if err != nil {
			return nil, err
		}
		if !t.Equal(exprType) {
			return nil, fmt.Errorf("type error at arg %d", i+1)
		}
	}

	return call.Return, nil
}

func isNumeric(typ Type) bool {
	return typ.Equal(Float) || typ.Equal(Integer)
}

func checkInfixExpr(expression *ast.InfixExpression, table SymbolTable) (Type, error) {
	leftType, leftErr := ExpressionType(expression.Left, table)
	if leftErr != nil {
		return nil, leftErr
	}

	rightType, rightErr := ExpressionType(expression.Right, table)
	if rightErr != nil {
		return nil, rightErr
	}

	// if the expression is an and/or, check for Boolean types
	if expression.Op == "&&" || expression.Op == "||" {
		if !leftType.Equal(Boolean) {
			return nil, fmt.Errorf("type error left operand")
		}
		if !rightType.Equal(Boolean) {
			return nil, fmt.Errorf("type error right operand")
		}
		return Boolean, nil
	}

	// if we don't have logical operators, we need to make sure that both
	// expressions are the same and are not Booleans
	if !isNumeric(leftType) {
		return nil, fmt.Errorf("type error: left operand not numerical")
	}
	if !isNumeric(rightType) {
		return nil, fmt.Errorf("type error: right operand not numerical")
	}

	if !leftType.Equal(rightType) {
		return nil, fmt.Errorf("type mismatch")
	}

	switch expression.Op {
	// comparison operators
	case "==", "!=", "<=", ">=", "<", ">":
		return Boolean, nil
	case "+", "-", "*", "/", "%":
		return leftType, nil
	default:
		panic("invalid infix operation")
	}
	return nil, nil
}

func checkPrefixExpr(expr *ast.PrefixExpression, table SymbolTable) (Type, error) {
	exprType, err := ExpressionType(expr.Expr, table)
	if err != nil {
		return nil, err
	}
	switch expr.Op {
	case "!":
		if !exprType.Equal(Boolean) {
			return nil, fmt.Errorf("type error, expected boolean on right hand side of '!'")
		}
		return Boolean, nil
	case "-":
		if !isNumeric(exprType) {
			return nil, fmt.Errorf("type error, expected numeric type on right hand side of '-'")
		}
		return exprType, nil
	default:
		panic("invalid prefix operation")
	}
}

func checkTupleRef(expr *ast.TupleRefExpression, table SymbolTable) (Type, error) {
	tup, err := ExpressionType(expr.Tuple, table)
	if err != nil {
		return nil, err
	}

	idx, ok := expr.Index.(*ast.IntExpression)
	if !ok {
		return nil, fmt.Errorf("expected integer literal received expression")
	}

	tupType, ok := tup.(*Tuple)
	if !ok {
		return nil, fmt.Errorf("tuple reference of non-tuple")
	}

	if idx.Val < 0 || int(idx.Val) > len(tupType.Types)-1 {
		return nil, fmt.Errorf("tuple index out of bounds")
	}

	return tupType.Types[idx.Val], nil
}

func checkArrayRef(expr *ast.ArrayRefExpression, table SymbolTable) (Type, error) {
	arr, err := ExpressionType(expr.Array, table)
	if err != nil {
		return nil, err
	}

	arrTyp, ok := arr.(*Array)
	if !ok {
		return nil, fmt.Errorf("array reference of non-array")
	}

	if arrTyp.Rank != len(expr.Indexes) {
		return nil, fmt.Errorf("array access of rank %d with %d indexes",
			arrTyp.Rank, len(expr.Indexes))
	}

	for _, idxExp := range expr.Indexes {
		idxTyp, err := ExpressionType(idxExp, table)
		if err != nil {
			return nil, err
		}
		if !idxTyp.Equal(Integer) {
			return nil, fmt.Errorf("non-integer index expression of array")
		}
	}

	return arrTyp.Inner, nil
}

func checkSumTransform(expr *ast.SumTransform, table SymbolTable) (Type, error) {
	cpy := table.Copy()
	for _, binding := range expr.OpBindings {
		if _, ok := cpy[binding.Variable]; ok {
			return nil, fmt.Errorf("illegal shadowing in sum expr, var: %s", binding.Variable)
		}
		bindType, err := ExpressionType(binding.Expr, cpy)
		if err != nil {
			return nil, err
		}

		if !bindType.Equal(Integer) {
			return nil, fmt.Errorf("bindArg expr initializer for %s returns non-integer",
				binding.Variable)
		}

		cpy[binding.Variable] = &Identifier{Type: Integer}
	}
	exprType, err := ExpressionType(expr.Expr, cpy)
	if err != nil {
		return nil, err
	}

	if !isNumeric(exprType) {
		return nil, fmt.Errorf("sum returns non-numeric expression")
	}

	return exprType, nil
}

func checkArrayTransform(expr *ast.ArrayTransform, table SymbolTable) (Type, error) {
	cpy := table.Copy()
	for _, binding := range expr.OpBindings {
		if _, ok := cpy[binding.Variable]; ok {
			return nil, fmt.Errorf("illegal shadowing in sum expr, var: %s", binding.Variable)
		}
		bindType, err := ExpressionType(binding.Expr, cpy)
		if err != nil {
			return nil, err
		}

		if !bindType.Equal(Integer) {
			return nil, fmt.Errorf("bindArg expr initializer for %s returns non-integer",
				binding.Variable)
		}

		cpy[binding.Variable] = &Identifier{Type: Integer}
	}
	exprType, err := ExpressionType(expr.Expr, cpy)
	if err != nil {
		return nil, err
	}

	if arr, ok := exprType.(*Array); ok {
		if arr.Rank != len(expr.OpBindings) {
			return nil, fmt.Errorf("return type of array expression must be of equal rank of number of bindings")
		}
		return exprType, nil
	}
	return nil, fmt.Errorf("return type of array expression must be array")
}

func checkTuple(expr *ast.TupleExpression, table SymbolTable) (Type, error) {
	tuple := &Tuple{
		Types: make([]Type, len(expr.Expressions)),
	}
	for i, expr := range expr.Expressions {
		typ, err := ExpressionType(expr, table)
		if err != nil {
			return nil, err
		}
		tuple.Types[i] = typ
	}

	return tuple, nil
}

func checkArray(expr *ast.ArrayExpression, table SymbolTable) (Type, error) {
	if len(expr.Expressions) == 0 {
		return &Array{Inner: Integer, Rank: 1}, nil
	}

	typ, err := ExpressionType(expr.Expressions[0], table)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(expr.Expressions); i++ {
		curTyp, err := ExpressionType(expr.Expressions[i], table)
		if err != nil {
			return nil, err
		}

		if !curTyp.Equal(typ) {
			return nil, fmt.Errorf("array literal has mixed types")
		}
	}

	return &Array{Inner: typ, Rank: 1}, nil
}

func ExpressionType(expression ast.Expression, table SymbolTable) (Type, error) {
	switch expr := expression.(type) {
	case *ast.BooleanExpression:
		return Boolean, nil
	case *ast.IntExpression:
		return Integer, nil
	case *ast.FloatExpression:
		return Float, nil
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
		panic("typechecking not implemented")

	}
	return nil, nil
}
