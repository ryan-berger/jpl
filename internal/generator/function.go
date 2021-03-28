package generator

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)


func calculateFunctionSize(stmts []ast.Statement) (int, frame) {
	size := 0
	f := make(frame)

	for _, cmd := range stmts {
		let, ok := cmd.(*ast.LetStatement)
		if !ok {
			continue
		}

		size += let.Expr.
			Typ().Size()

		// TODO: Later, we will need to do a switch on this to handle the rank of the variable
		ident := let.LValue.(*ast.VariableArgument).Variable
		f[ident] = size
	}

	if extra := size % 16; extra != 0 { // pad the stack frame
		size += extra
	}

	if size == 0 { // if the stack size is 0, it is not a multiple of 16
		size = 16
	}

	return size, f
}

func (g *generator) genFunction(f *ast.Function) {
	size, frame := calculateFunctionSize(f.Statements)
	g.buf.WriteString(fmt.Sprintf("_%s:\n", f.Var))
	for _, stmt := range f.Statements {
		g.genStatement(stmt, frame, size)
	}
	g.buf.WriteByte('\n')
}
