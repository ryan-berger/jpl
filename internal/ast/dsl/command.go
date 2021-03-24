package dsl

import "github.com/ryan-berger/jpl/internal/ast"

func Print(str string) *ast.Print {
	return &ast.Print{Str: str}
}

func Show(expr ast.Expression) *ast.Show {
	return &ast.Show{Expr: expr}
}
