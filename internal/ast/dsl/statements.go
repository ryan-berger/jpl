package dsl

import "github.com/ryan-berger/jpl/internal/ast"

func Return(exp ast.Expression) *ast.ReturnStatement {
	return &ast.ReturnStatement{Expr: exp}
}

func Assert(exp ast.Expression, message string) *ast.AssertStatement {
	return &ast.AssertStatement{
		Expr: exp,
		Message: message,
	}
}

func Attribute(annotation string) *ast.AttributeStatement {
	return &ast.AttributeStatement{Annotation: annotation}
}
