package generator

import (
	"bytes"
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

const prologue = `section .text

_main:
`

func calculateProgramSize(program ast.Program) (frame, int) {
	size := 0
	f := make(frame)
	for _, cmd := range program {
		let, ok := cmd.(*ast.LetStatement)
		if !ok {
			continue
		}

		size += let.Expr.
			Typ().Size()
		ident := let.LValue.(*ast.VariableArgument).Variable
		f[ident] = size
	}

	if extra := size % 16; extra != 0 {
		size += extra
	}

	return f, size
}

const fnPrologue = `push rbp
mov rbp, rsp
sub rsp %d`

func textSection(program ast.Program, mapper constantMapper) string {
	buf := bytes.NewBufferString(prologue)

	_, programSize := calculateProgramSize(program)

	buf.WriteString(fmt.Sprintf(fnPrologue, programSize))
	fmt.Println(buf.String())


	return buf.String()
}
