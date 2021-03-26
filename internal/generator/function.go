package generator

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)



func (g *generator) genFunction(f *ast.Function) {
	g.buf.WriteString(fmt.Sprintf("_%s:\n", f.Var))
	for _, stmt := range f.Statements {
		g.genStatement(stmt)
	}
	g.buf.WriteByte('\n')
}
