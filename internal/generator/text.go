package generator

import (
	"bytes"
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)

const textProlouge = `section .text

main:
_main:

`

func (g *generator) calculateProgramSize() {
	size := 0
	f := make(frame)

	for _, cmd := range g.program {
		if read, ok := cmd.(*ast.Read); ok {
			size += types.Pict.Size()

			ident := read.Argument.(*ast.VariableArgument).Variable
			f[ident] = size
		}

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

	if size == 0 {
		size = 16
	}

	g.frame = f
	g.size = size

	g.buf.WriteString(fmt.Sprintf(fnPrologue, size))
}

const fnPrologue = `push rbp
mov rbp, rsp
sub rsp, %d

`

func (g *generator) textSection() {
	g.buf.WriteString(textProlouge)

	g.calculateProgramSize()

	for _, cmd := range g.program {
		g.genCommand(cmd)
	}
}

type generator struct {
	program ast.Program
	table   *symbol.Table
	mapper  constantMapper

	frame frame
	size  int

	curLabel int

	buf *bytes.Buffer
}

func (g *generator) generate() string {
	g.buf.WriteString(filePrologue)

	dataString, mapper := dataSection(g.program)
	g.buf.WriteString(dataString)

	g.mapper = mapper
	g.textSection()

	return g.buf.String()
}

func (g *generator) newLabel() string {
	lbl := fmt.Sprintf("lbl%d", g.curLabel)
	g.curLabel++
	return lbl
}
