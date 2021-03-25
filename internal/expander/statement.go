package expander

import (
	"github.com/ryan-berger/jpl/internal/ast"
)

func expandStatement(statement ast.Statement, next nexter) []ast.Statement {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		return expandLet(stmt, next)
	case *ast.AssertStatement:
		return expandAssert(stmt, next)
	case *ast.ReturnStatement:
		return expandReturn(stmt, next)
	}
	return nil
}

func expandLet(ls *ast.LetStatement, next nexter) []ast.Statement {
	switch ls.LValue.(type) {
	case *ast.VariableArgument:
		exp, stmts := expandExpression(ls.Expr, next)
		ls.Expr = exp
		return append(stmts, ls)
	}
	return nil
}

func expandAssert(a *ast.AssertStatement, next nexter) []ast.Statement {
	ref, stmts := expansionAndLet(a.Expr, next)
	a.Expr = ref

	return append(stmts, a)
}

func expandReturn(r *ast.ReturnStatement, next nexter) []ast.Statement {
	ref, stmts := expansionAndLet(r.Expr, next)
	r.Expr = ref
	return append(stmts, r)
}