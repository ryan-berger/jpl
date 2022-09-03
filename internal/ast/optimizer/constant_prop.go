package optimizer

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

func isPropable(expr ast.Expression) bool {
	switch expr.(type) {
	case *ast.BooleanExpression, *ast.IntExpression, *ast.FloatExpression:
		return true
	}
	return false
}

func replace(expression ast.Expression, ident string, with ast.Expression) ast.Expression {
	switch expr := expression.(type) {
	case *ast.IdentifierExpression:
		fmt.Printf("replacing: %v with %v\n", ident, with)
		if expr.Identifier == ident {
			return with
		}
	case *ast.SumTransform:
		for _, bin := range expr.OpBindings {
			bin.Expr = replace(bin.Expr, ident, with)
		}
		expr.Expr = replace(expr.Expr, ident, with)
	case *ast.ArrayTransform:
		for _, bin := range expr.OpBindings {
			bin.Expr = replace(bin.Expr, ident, with)
		}
		expr.Expr = replace(expr.Expr, ident, with)
	case *ast.InfixExpression:
		expr.Left = replace(expr.Left, ident, with)
		expr.Right = replace(expr.Right, ident, with)
	case *ast.PrefixExpression:
		expr.Expr = replace(expr.Expr, ident, with)
	case *ast.CallExpression:
		for i, args := range expr.Arguments {
			expr.Arguments[i] = replace(args, ident, with)
		}
	case *ast.ArrayExpression:
		for i, args := range expr.Expressions {
			expr.Expressions[i] = replace(args, ident, with)
		}
	case *ast.IfExpression:
		expr.Condition = replace(expr.Condition, ident, with)
		expr.Consequence = replace(expr.Consequence, ident, with)
		expr.Otherwise = replace(expr.Otherwise, ident, with)
	}
	return expression
}

func propCmd(n ast.Node, ident string, with ast.Expression) {
	switch node := n.(type) {
	case ast.Statement:
		propStmt(node, ident, with)
	case *ast.Show:
		node.Expr = replace(node.Expr, ident, with)
	case *ast.Time:
		propCmd(node.Command, ident, with)
	}
}

func propStmt(statement ast.Statement, ident string, with ast.Expression) {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		stmt.Expr = replace(stmt.Expr, ident, with)
	case *ast.AssertStatement:
		stmt.Expr = replace(stmt.Expr, ident, with)
	case *ast.ReturnStatement:
		stmt.Expr = replace(stmt.Expr, ident, with)
	}
}

func checkAndProp(n ast.Node, use *defUse) {
	switch stmt := n.(type) {
	case *ast.LetStatement:
		if isPropable(stmt.Expr) {
			variable := stmt.LValue.(*ast.Variable).Variable
			for _, u := range use.getUses(variable) {
				propCmd(u, variable, stmt.Expr)
				use.clearUse(n)
			}
			delete(use.graph, variable)
		}
	case *ast.Function:
		for _, s := range stmt.Statements {
			checkAndProp(s, use)
		}
	}
}

func ConstantProp(program ast.Program) ast.Program {
	du := buildDefUse(program)

	for _, cmd := range program {
		checkAndProp(cmd, du)
	}

	return program
}
