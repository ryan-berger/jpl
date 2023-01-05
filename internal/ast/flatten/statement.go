package flatten

import (
	"github.com/ryan-berger/jpl/internal/ast"
)

func flattenStatement(statement ast.Statement, next nexter) []ast.Statement {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		return flattenLet(stmt, next)
	case *ast.AssertStatement:
		return flattenAssert(stmt, next)
	case *ast.ReturnStatement:
		return flattenReturn(stmt, next)
	}
	return nil
}

func flattenLet(ls *ast.LetStatement, next nexter) []ast.Statement {
	exp, stmts := flattenExpression(ls.Expr, next)
	ls.Expr = exp
	return append(stmts, ls)
}

func flattenAssert(a *ast.AssertStatement, next nexter) []ast.Statement {
	ref, stmts := expansionAndLet(a.Expr, next)
	a.Expr = ref

	return append(stmts, a)
}

func flattenReturn(r *ast.ReturnStatement, next nexter) []ast.Statement {
	ref, stmts := expansionAndLet(r.Expr, next)
	r.Expr = ref
	return append(stmts, r)
}
