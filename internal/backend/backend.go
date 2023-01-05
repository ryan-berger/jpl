package backend

import (
	"io"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/symbol"
)

type Generator func(program ast.Program, st *symbol.Table, writer io.Writer)
