package optimizer

import (
	"errors"
	"math"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

func isConstant(exp ast.Expression) bool {
	switch exp.(type) {
	case *ast.IntExpression, *ast.FloatExpression:
		return true
	}
	return false
}


var DivideByZero = errors.New("divide by zero")

func foldInteger(l, r ast.Expression, op string) (ast.Expression, error) {
	lInt := l.(*ast.IntExpression).Val
	rInt := r.(*ast.IntExpression).Val
	if rInt == 0 && op == "/" || op == "%" {
		return nil, DivideByZero
	}

	var exp ast.Expression
	switch op {
	case "+":
		exp = dsl.Int(lInt + rInt)
	case "-":
		exp = dsl.Int(lInt - rInt)
	case "/":
		exp = dsl.Int(lInt / rInt)
	case "*":
		exp = dsl.Int(lInt * rInt)
	case "%":
		exp = dsl.Int(lInt % rInt)
	case "<":
		exp = dsl.Bool(lInt < rInt)
	case ">":
		exp = dsl.Bool(lInt > rInt)
	case "<=":
		exp = dsl.Bool(lInt <= rInt)
	case ">=":
		exp = dsl.Bool(lInt >= rInt)
	case "==":
		exp = dsl.Bool(lInt == rInt)
	case "!=":
		exp = dsl.Bool(lInt != rInt)
	default:
		panic("error, not implemented")
	}
	return exp, nil
}

func foldFloat(l, r ast.Expression, op string) ast.Expression {
	lFloat := l.(*ast.FloatExpression).Val
	rFloat := r.(*ast.FloatExpression).Val
	switch op {
	case "+":
		return dsl.Float(lFloat + rFloat)
	case "-":
		return dsl.Float(lFloat - rFloat)
	case "/":
		return dsl.Float(lFloat / rFloat)
	case "%":
		return dsl.Float(math.Mod(lFloat, rFloat))
	case "*":
		return dsl.Float(lFloat * rFloat)
	case "<":
		return dsl.Bool(lFloat < rFloat)
	case ">":
		return dsl.Bool(lFloat > rFloat)
	case "<=":
		return dsl.Bool(lFloat <= rFloat)
	case ">=":
		return dsl.Bool(lFloat >= rFloat)
	case "==":
		return dsl.Bool(lFloat == rFloat)
	case "!=":
		return dsl.Bool(lFloat != rFloat)
	default:
		panic("error, not implemented")
	}
}

func constantFold(exp ast.Expression) (ast.Expression, error) {
	switch expr := exp.(type) {
	case *ast.IntExpression, *ast.FloatExpression, *ast.BooleanExpression:
		return expr, nil
	case *ast.ArrayExpression:
		for i, e := range expr.Expressions {
			f, err := constantFold(e)
			if err != nil {
				return nil, err
			}
			expr.Expressions[i] = f
		}
		return expr, nil
	case *ast.TupleExpression:
		for i, e := range expr.Expressions {
			f, err := constantFold(e)
			if err != nil {
				return nil, err
			}
			expr.Expressions[i] = f
		}
		return expr, nil
	case *ast.CallExpression:
		for i, e := range expr.Arguments {
			f, err := constantFold(e)
			if err != nil {
				return nil, err
			}
			expr.Arguments[i] = f
		}
		return expr, nil
	case *ast.ArrayTransform:
		if len(expr.OpBindings) == 0 {
			return constantFold(expr.Expr)
		}

		for i, bind := range expr.OpBindings {
			exp, err := constantFold(bind.Expr)
			if err != nil {
				return nil, err
			}
			expr.OpBindings[i].Expr = exp
		}

		exp, err := constantFold(expr.Expr)
		if err != nil {
			return nil, err
		}

		expr.Expr = exp
		return expr, nil
	case *ast.SumTransform:
		if len(expr.OpBindings) == 0 {
			return constantFold(expr.Expr)
		}

		for i, bind := range expr.OpBindings {
			exp, err := constantFold(bind.Expr)
			if err != nil {
				return nil, err
			}
			expr.OpBindings[i].Expr = exp
		}

		exp, err := constantFold(expr.Expr)
		if err != nil {
			return nil, err
		}

		expr.Expr = exp
		return expr, nil
	case *ast.IfExpression:
		var err error
		expr.Condition, err = constantFold(expr.Condition)
		if err != nil {
			return nil, err
		}

		expr.Consequence, err = constantFold(expr.Consequence)
		if err != nil {
			return nil, err
		}

		expr.Otherwise, err = constantFold(expr.Otherwise)
		if err != nil {
			return nil, err
		}

		return expr, nil
	case *ast.InfixExpression:
		lExp, err := constantFold(expr.Left)
		if err != nil {
			return nil, err
		}
		rExp, err := constantFold(expr.Right)
		if err != nil {
			return nil, err
		}
		if isConstant(lExp) && isConstant(rExp) {
			switch lExp.(type) {
			case *ast.FloatExpression:
				return foldFloat(lExp, rExp, expr.Op), nil
			case *ast.IntExpression:
				return foldInteger(lExp, rExp, expr.Op)
			default:
				panic("no other type possible")
			}
		}
		return dsl.Infix(expr.Op, lExp, rExp), nil
	case *ast.PrefixExpression:
		r, err := constantFold(expr.Expr)
		if err != nil {
			return nil, err
		}

		if expr.Op == "-" {
			switch exp := expr.Expr.(type) {
			case *ast.IntExpression:
				return dsl.Int(-exp.Val), nil
			case *ast.FloatExpression:
				return dsl.Float(-exp.Val), nil
			default:
				panic("oops")
			}
		}
		return dsl.Prefix(expr.Op, r), nil
	default:
		return exp, nil
	}
}

func foldStmt(cmd ast.Statement) (ast.Statement, error) {
	switch c := cmd.(type) {
	case *ast.LetStatement:
		exp, err := constantFold(c.Expr)
		if err == DivideByZero {
			return dsl.Assert(dsl.Bool(false), divideByZero), err
		}
		c.Expr = exp
		return c, nil
	case *ast.AssertStatement:
		exp, err := constantFold(c.Expr)
		if err == DivideByZero {
			return dsl.Assert(dsl.Bool(false), divideByZero), err
		}
		c.Expr = exp
		return c, nil
	case *ast.ReturnStatement:
		exp, err := constantFold(c.Expr)
		if err == DivideByZero {
			return dsl.Assert(dsl.Bool(false), divideByZero), err
		}
		c.Expr = exp
		return c, nil
	default:
		panic("oops")
	}
}

func foldCommand(cmd ast.Command) (ast.Command, error) {
	switch c := cmd.(type) {
	case ast.Statement:
		return foldStmt(c)
	case *ast.Show:
		exp, err := constantFold(c.Expr)
		if err == DivideByZero {
			return dsl.Assert(dsl.Bool(false), divideByZero), err
		}
		c.Expr = exp
		return c, nil
	case *ast.Time:
		newCmd, err := foldCommand(c.Command)
		if err != nil {
			return newCmd, err
		}
		c.Command = newCmd
		return c, nil
	case *ast.Function:
		var stmts []ast.Statement
		for _, stmt := range c.Statements {
			newStmt, err := foldStmt(stmt)
			stmts = append(stmts, newStmt)
			if err != nil {
				break
			}
		}
		c.Statements = stmts
		return c, nil
	}
	panic("yeet")
}

const divideByZero = "Division by zero error"

func ConstantFold(p ast.Program) ast.Program {
	var next ast.Program
	for _, cmd := range p {
		newCmd, err := foldCommand(cmd)
		next = append(next, newCmd)
		if err != nil {
			break
		}
	}
	return next
}
