package optimizer

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

func sumExpr(expression ast.Expression) ast.Expression {
	switch exp := expression.(type) {
	case *ast.InfixExpression:
		exp.Left = sumExpr(exp.Left)
		exp.Right = sumExpr(exp.Right)
	case *ast.PrefixExpression:
		exp.Expr = sumExpr(exp.Expr)
	case *ast.CallExpression:
		for i := 0; i < len(exp.Arguments); i++ {
			exp.Arguments[i] = sumExpr(exp.Arguments[i])
		}
	case *ast.IfExpression:
		exp.Condition = sumExpr(exp.Condition)
		exp.Consequence = sumExpr(exp.Consequence)
		exp.Otherwise = sumExpr(exp.Otherwise)
	case *ast.ArrayTransform:
		for _, s := range exp.OpBindings {
			s.Expr = sumExpr(s.Expr)
		}
		exp.Expr = sumExpr(exp.Expr)
	case *ast.SumTransform:
		var product int64 = 1
		for _, s := range exp.OpBindings {
			s.Expr = sumExpr(s.Expr)
			v, ok := s.Expr.(*ast.IntExpression)
			if !ok {
				return exp
			}
			product *= v.Val
		}
		rh, ok := exp.Expr.(*ast.IntExpression)
		if !ok {
			return exp
		}
		return dsl.Int(product * rh.Val)
	}
	return expression
}

func sumStmt(statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		stmt.Expr = sumExpr(stmt.Expr)
	case *ast.AssertStatement:
		stmt.Expr = sumExpr(stmt.Expr)
	case *ast.ReturnStatement:
		stmt.Expr = sumExpr(stmt.Expr)
	}
}

func sumCmd(command ast.Command) {
	switch cmd := command.(type) {
	case ast.Statement:
		sumStmt(cmd)
	case *ast.Time:
		sumCmd(cmd.Command)
	case *ast.Function:
		for _, c := range cmd.Statements {
			sumCmd(c)
		}
	case *ast.Show:
		cmd.Expr = sumExpr(cmd.Expr)
	case *ast.Write:
		cmd.Expr = sumExpr(cmd.Expr)
	}
}

func Peephole(p ast.Program) {
	for _, c := range p {
		sumCmd(c)
	}
}
