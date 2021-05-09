package optimizer

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

func sumExpr(expression ast.Expression) (ast.Expression, error) {
	var err error
	switch exp := expression.(type) {
	case *ast.InfixExpression:
		exp.Left, err = sumExpr(exp.Left)
		if err != nil {
			return nil, err
		}
		exp.Right, err = sumExpr(exp.Right)
		if err != nil {
			return nil, err
		}
	case *ast.PrefixExpression:
		exp.Expr, err = sumExpr(exp.Expr)
		if err != nil {
			return nil, err
		}
	case *ast.CallExpression:
		for i := 0; i < len(exp.Arguments); i++ {
			exp.Arguments[i], err = sumExpr(exp.Arguments[i])
			if err != nil {
				return nil, err
			}
		}
	case *ast.IfExpression:
		exp.Condition, err = sumExpr(exp.Condition)
		if err != nil {
			return nil, err
		}
		exp.Consequence, err = sumExpr(exp.Consequence)
		if err != nil {
			return nil, err
		}
		exp.Otherwise, err = sumExpr(exp.Otherwise)
		if err != nil {
			return nil, err
		}
	case *ast.ArrayTransform:
		for _, s := range exp.OpBindings {
			s.Expr, err = sumExpr(s.Expr)
			if err != nil {
				return nil, err
			}
		}
		exp.Expr, err = sumExpr(exp.Expr)
		if err != nil {
			return nil, err
		}
	case *ast.SumTransform:
		var product int64 = 1
		for _, s := range exp.OpBindings {
			s.Expr, err = sumExpr(s.Expr)
			v, ok := s.Expr.(*ast.IntExpression)
			if !ok {
				return exp, nil
			}
			product *= v.Val
		}
		exp.Expr, err = constantFold(exp.Expr)
		rh, ok := exp.Expr.(*ast.IntExpression)
		if !ok {
			return exp, nil
		}
		return dsl.Int(product * rh.Val), nil
	}
	return expression, nil
}

func sumStmt(statement ast.Statement) ast.Statement {
	var err error
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		stmt.Expr, err = sumExpr(stmt.Expr)
	case *ast.AssertStatement:
		stmt.Expr, err = sumExpr(stmt.Expr)
	case *ast.ReturnStatement:
		stmt.Expr, err = sumExpr(stmt.Expr)
	}

	if err != nil {
		return dsl.Assert(dsl.Bool(false), divideByZero)
	}
	return statement
}

func sumCmd(command ast.Command) ast.Command {
	var err error
	switch cmd := command.(type) {
	case ast.Statement:
		return sumStmt(cmd)
	case *ast.Time:
		return sumCmd(cmd.Command)
	case *ast.Function:
		for i, c := range cmd.Statements {
			cmd.Statements[i] = sumStmt(c)
		}
	case *ast.Show:
		cmd.Expr, err = sumExpr(cmd.Expr)
		if err != nil {
			return dsl.Assert(dsl.Bool(false), divideByZero)
		}
	case *ast.Write:
		cmd.Expr, err = sumExpr(cmd.Expr)
		if err != nil {
			return dsl.Assert(dsl.Bool(false), divideByZero)
		}
	}
	return command
}

func Peephole(p ast.Program) ast.Program {
	for i, c := range p {
		p[i] = sumCmd(c)
	}
	return p
}
