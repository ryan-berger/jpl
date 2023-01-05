package dsl

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/types"
)

func Fn(name string, bindings []ast.Binding, ret types.Type, stmts []ast.Statement) *ast.Function {
	return &ast.Function{
		Var:        name,
		Bindings:   bindings,
		ReturnType: ret,
		Statements: stmts,
	}
}
