package flatten

import (
	"github.com/ryan-berger/jpl/internal/ast"
)

func toCommands(stmts []ast.Statement) []ast.Command {
	cmds := make([]ast.Command, len(stmts))
	for i, stmt := range stmts {
		cmds[i] = stmt
	}
	return cmds
}
