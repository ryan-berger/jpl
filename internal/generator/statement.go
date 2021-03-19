package generator

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/types"
)

func ident(let *ast.LetStatement) string {
	arg, ok := let.LValue.(*ast.VariableArgument)
	if !ok {
		panic("must call ident() on VariableArgument")
	}
	return arg.Variable
}

func genCommand(command ast.Command, f frame) {

}

func genStatement(
	statement ast.Statement,
	consts map[ast.Statement]string,
	f frame) string {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		return genLetStatement(stmt, consts, f)
	}
	return ""
}

const letConstant = `mov rbx, [rel %s]
mov [rbp - %d], rbx`

const moveVar = `mov rbx [rbp - %d]
mov [rbp - %d], rbp`

func genLetStatement(
	let *ast.LetStatement,
	consts map[ast.Statement]string,
	f frame) string {

	switch exp := let.Expr.(type) {
	case *ast.FloatExpression, *ast.IntExpression:
		constName := consts[let]
		loc := f[ident(let)]
		return fmt.Sprintf(letConstant, constName, loc)
	case *ast.IdentifierExpression:
		typ := exp.Type
		switch {
		case typ.Equal(types.Float) ||
			typ.Equal(types.Integer):
			locL := f[ident(let)]
			locR := f[exp.Identifier]

			return fmt.Sprintf(moveVar,locL, locR)
		}

	}
	return ""
}
