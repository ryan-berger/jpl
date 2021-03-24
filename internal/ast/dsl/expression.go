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

func Ident(val string) *ast.IdentifierExpression {
	return &ast.IdentifierExpression{
		Identifier: val,
	}
}
