package optimizer

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

func isConstant(exp ast.Expression) bool {
	switch exp.(type) {
	case *ast.IntExpression, *ast.FloatExpression, *ast.BooleanExpression:
		return true
	}
	return false
}

func foldBoolExpr(l, r ast.Expression, op string) ast.Expression {
	lBool := l.(*ast.BooleanExpression).Val
	rBool := r.(*ast.BooleanExpression).Val

	switch op {
	case "&&":
		return dsl.Bool(lBool && rBool)
	case "||":
		return dsl.Bool(lBool || rBool)
	default:
		panic("error, not implemented")
	}
}

func foldInteger(l, r ast.Expression, op string) ast.Expression {
	lInt := l.(*ast.IntExpression).Val
	rInt := r.(*ast.IntExpression).Val
	if rInt == 0 { // TODO: actually handle this
		panic("divide by zero")
	}
	switch op {
	case "+":
		return dsl.Int(lInt + rInt)
	case "-":
		return dsl.Int(lInt - rInt)
	case "/":
		return dsl.Int(lInt / rInt)
	case "*":
		return dsl.Int(lInt * rInt)
	case "<":
		return dsl.Bool(lInt < rInt)
	case ">":
		return dsl.Bool(lInt > rInt)
	case "<=":
		return dsl.Bool(lInt <= rInt)
	case ">=":
		return dsl.Bool(lInt >= rInt)
	case "==":
		return dsl.Bool(lInt == rInt)
	case "!=":
		return dsl.Bool(lInt != rInt)
	default:
		panic("error, not implemented")
	}
}

func foldFloat(l, r ast.Expression, op string) ast.Expression {
	lInt := l.(*ast.FloatExpression).Val
	rInt := r.(*ast.FloatExpression).Val
	switch op {
	case "+":
		return dsl.Float(lInt + rInt)
	case "-":
		return dsl.Float(lInt - rInt)
	case "/":
		return dsl.Float(lInt / rInt)
	case "*":
		return dsl.Float(lInt * rInt)
	case "<":
		return dsl.Bool(lInt < rInt)
	case ">":
		return dsl.Bool(lInt > rInt)
	case "<=":
		return dsl.Bool(lInt <= rInt)
	case ">=":
		return dsl.Bool(lInt >= rInt)
	case "==":
		return dsl.Bool(lInt == rInt)
	case "!=":
		return dsl.Bool(lInt != rInt)
	default:
		panic("error, not implemented")
	}
}


func constantFold(exp ast.Expression) ast.Expression {
	switch expr := exp.(type) {
	case *ast.IntExpression, *ast.FloatExpression, *ast.BooleanExpression:
		return expr
	case *ast.SumTransform:
		allConstant := true

		for i, bind := range expr.OpBindings {
			exp := constantFold(bind.Expr)
			allConstant = allConstant && isConstant(exp)
			expr.OpBindings[i].Expr = exp
		}

		return expr
	case *ast.InfixExpression:
		lExp := constantFold(expr.Left)
		rExp := constantFold(expr.Right)

		if isConstant(lExp) && isConstant(rExp) {
			switch lExp.(type) {
			case *ast.FloatExpression:
				return foldFloat(lExp, rExp, expr.Op)
			case *ast.IntExpression:
				return foldInteger(lExp, rExp, expr.Op)
			case *ast.BooleanExpression:
				return foldBoolExpr(lExp, rExp, expr.Op)
			}
		}
		return dsl.Infix(expr.Op, lExp, rExp)
	case *ast.IfExpression:
		cons := constantFold(expr.Consequence)
		other := constantFold(expr.Otherwise)
		cond := constantFold(expr)
		if isConstant(cond) {
			if cond.(*ast.BooleanExpression).Val {
				return cons
			}
			return other
		}
		return dsl.If(cond, cons, other)
	case *ast.PrefixExpression:
		return expr
	default:
		return exp
	}
}

