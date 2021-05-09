package nasm

import (
	"bytes"
	"io"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
)

const filePrologue = `global main
global _main

extern _fail_assertion
extern _sub_floats
extern _sub_ints
extern _has_size
extern _sepia
extern _blur
extern _resize
extern _crop
extern _read_image
extern _print
extern _write_image
extern _show
extern _fail_assertion
extern _get_time

`



func Generate(program ast.Program, table *symbol.Table, out io.Writer) {
	g := generator{
		program: program,
		table:   table,
		buf:     bytes.NewBuffer([]byte{}),
	}

	out.Write([]byte(g.generate()))
}
