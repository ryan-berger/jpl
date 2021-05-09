package dsl

import "github.com/ryan-berger/jpl/internal/ast"

func LTuple(vals ...ast.LValue) *ast.LTuple {
	return &ast.LTuple{Args: vals}
}

func LArray(ident string, dims ...string) *ast.VariableArr {
	return &ast.VariableArr{
		Variable:  ident,
		Variables: dims,
	}
}

func LIdent(ident string) *ast.Variable {
	return &ast.Variable{
		Variable: ident,
	}
}

func Let(lval ast.LValue, expr ast.Expression) *ast.LetStatement {
	return &ast.LetStatement{
		LValue: lval,
		Expr:   expr,
	}
}
