package expander

import "github.com/ryan-berger/jpl/internal/ast"

// expansionAndLet takes in an expression and if it is an ast.CallExpression it will convert
// it to a let. This method should only be used when an ast.CallExpression is not allowed in an
// expansion and allows for expandExpression to stay generic
func expansionAndLet(
	expr ast.Expression,
	next nexter,
) (ast.Expression, []ast.Statement) {
	if _, ok := expr.(*ast.IdentifierExpression); ok {
		return expr, nil
	}

	newExp, stmts := expandExpression(expr, next)
	l := let(next(), newExp)
	stmts = append(stmts, l)

	return refExpr(ident(l.LValue)), stmts
}

func expandExpression(expression ast.Expression, next nexter) (ast.Expression, []ast.Statement) {
	switch expr := expression.(type) {
	case *ast.IdentifierExpression:
		return expr, nil
	case *ast.IntExpression, *ast.FloatExpression:
		return expr, []ast.Statement{}
	case *ast.CallExpression:
		var stmts []ast.Statement
		for i, arg := range expr.Arguments {
			newExp, deps := expansionAndLet(arg, next) // guarantee that we are a let

			expr.Arguments[i] = newExp
			stmts = append(stmts, deps...)
		}
		return expr, stmts
	}
	return nil, nil
}
