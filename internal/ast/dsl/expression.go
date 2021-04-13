package dsl

import "github.com/ryan-berger/jpl/internal/ast"

func Call(name string, args ...ast.Expression) *ast.CallExpression {
	return &ast.CallExpression{
		Identifier: name,
		Arguments:  args,
	}
}

func Int(val int64) *ast.IntExpression {
	return &ast.IntExpression{Val: val}
}

func Float(val float64) *ast.FloatExpression {
	return &ast.FloatExpression{Val: val}
}

func Bool(val bool) *ast.BooleanExpression {
	return &ast.BooleanExpression{Val: val}
}

func Ident(val string) *ast.IdentifierExpression {
	return &ast.IdentifierExpression{
		Identifier: val,
	}
}

func Tuple(exps ...ast.Expression) *ast.TupleExpression {
	return &ast.TupleExpression{Expressions: exps}
}

func Infix(op string, l, r ast.Expression) *ast.InfixExpression {
	return &ast.InfixExpression{
		Left:  l,
		Right: r,
		Op:    op,
	}
}

func Prefix(op string, r ast.Expression) *ast.PrefixExpression {
	return &ast.PrefixExpression{
		Op:   op,
		Expr: r,
	}
}

func If(cond, otherwise, consequence ast.Expression) *ast.IfExpression {
	return &ast.IfExpression{
		Condition:   cond,
		Consequence: consequence,
		Otherwise:   otherwise,
	}
}
