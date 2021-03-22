package expander

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

func let(counter int, expr ast.Expression) *ast.LetStatement {
	return &ast.LetStatement{
		LValue: &ast.VariableArgument{Variable: fmt.Sprintf("t.%d", counter)},
		Expr:   expr,
	}
}

func constInt(val int64) *ast.IntExpression {
	return &ast.IntExpression{Val: val}
}

func returnStmt(expr ast.Expression) *ast.ReturnStatement {
	return &ast.ReturnStatement{Expr: expr}
}

func ident(value ast.LValue) string {
	arg, ok := value.(*ast.VariableArgument)
	if !ok {
		panic("must call ident() on VariableArgument")
	}
	return arg.Variable
}

func refExpr(ref string) *ast.IdentifierExpression {
	return &ast.IdentifierExpression{Identifier: ref}
}

func print(str string) *ast.Print {
	return &ast.Print{
		Str: fmt.Sprintf(`"%s"`, str),
	}
}

func show(expr ast.Expression) ast.Command {
	return &ast.Show{
		Expr: expr,
	}
}

func functionCall(name string, args ...ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Identifier: name,
		Arguments:  args,
	}
}

func toCommands(stmts []ast.Statement) []ast.Command {
	cmds := make([]ast.Command, len(stmts))
	for i, stmt := range stmts {
		cmds[i] = stmt
	}
	return cmds
}

