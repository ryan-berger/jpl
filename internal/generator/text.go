package generator

import (
	"bytes"
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

const prologue = `section .text

_main:
`

type frame struct {
}

func calculateProgramSize(program ast.Program) int {
	size := 0
	for _, cmd := range program {
		let, ok := cmd.(*ast.LetStatement)
		if !ok {
			continue
		}
		size += let.Expr.
			Typ().Size()
	}
	return size
}

func textSection(program ast.Program, mapper constantMapper) string {
	buf := bytes.NewBufferString(prologue)

	programSize := calculateProgramSize(program)
	fmt.Println(programSize)

	return buf.String()
}
