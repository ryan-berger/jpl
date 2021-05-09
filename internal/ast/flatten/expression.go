package flatten

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

// expansionAndLet takes in an expression and if it is an ast.CallExpression it will convert
// it to a let. This method should only be used when an ast.CallExpression is not allowed in an
// expansion and allows for flattenExpression to stay generic
func expansionAndLet(
	expr ast.Expression,
	next nexter,
) (*ast.IdentifierExpression, []ast.Statement) {
	if e, ok := expr.(*ast.IdentifierExpression); ok {
		return e, nil
	}

	newExp, stmts := flattenExpression(expr, next)
	name := next()
	l := dsl.Let(
		dsl.LIdent(name), newExp) // let ident = newExp

	stmts = append(stmts, l)
	return dsl.Ident(name), stmts
}

func flattenInfixExpression(expr ast.Expression, next nexter) (ast.Expression, []ast.Statement) {
	switch exp := expr.(type) {
	case *ast.FloatExpression, *ast.IntExpression:
		return expansionAndLet(expr, next)
	case *ast.InfixExpression:
		if exp.Op == "&&" || exp.Op == "||" { // TODO: generate function calls for these
			return exp, nil
		}

		lExp, lStatements := expansionAndLet(exp.Left, next)
		rExp, rStatements := expansionAndLet(exp.Right, next)
		var stmts []ast.Statement
		stmts = append(stmts, lStatements...)
		stmts = append(stmts, rStatements...)

		return &ast.InfixExpression{Left: lExp, Right: rExp, Op: exp.Op}, stmts
	default:
		return flattenExpression(expr, next)
	}
}


func flattenExpression(expression ast.Expression, next nexter) (ast.Expression, []ast.Statement) {
	switch expr := expression.(type) {
	case *ast.IdentifierExpression:
		return expr, nil
	case *ast.IntExpression, *ast.FloatExpression, *ast.BooleanExpression:
		return expr, nil
	case *ast.IfExpression:
		return expr, nil
	case *ast.InfixExpression:
		return flattenInfixExpression(expr, next)
	case *ast.TupleExpression:
		return expr, nil
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
