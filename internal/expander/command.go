package expander

import (
	"github.com/ryan-berger/jpl/internal/ast"
)

func expandCommand(command ast.Command, next nexter) []ast.Command {
	switch cmd := command.(type) {
	case *ast.Read:
		return []ast.Command{cmd}
	case *ast.Print:
		return []ast.Command{cmd}
	case *ast.Write:
		return expandWrite(cmd, next)
	case *ast.Time:
		return expandTime(cmd, next)
	case *ast.Show:
		return expandShow(cmd, next)
	case ast.Statement:
		return toCommands(expandStatement(cmd, next))
	default:
		panic("oops, type not supported")
	}
	return nil
}

func expandShow(sh *ast.Show, next nexter) []ast.Command {
	ref, cmds := expansionAndLet(sh.Expr, next)

	sh.Expr = ref

	var combined []ast.Command
	combined = append(combined, toCommands(cmds)...)
	combined = append(combined, sh)
	return combined
}

func expandTime(time *ast.Time, next nexter) []ast.Command {
	start := let(next(), functionCall("get_time"))
	cmds := expandCommand(time.Command, next)
	end := let(next(), functionCall("get_time"))
	sub := let(next(),
		functionCall("sub_floats",
			refExpr(ident(start.LValue)),
			refExpr(ident(end.LValue))))
	p := print("time: ")
	s := show(refExpr(ident(sub.LValue)))

	commands := []ast.Command{start}
	commands = append(commands, cmds...)
	commands = append(commands, sub, p, s)
	return cmds
}

func expandWrite(w *ast.Write, next nexter) []ast.Command {
	ref, cmds := expansionAndLet(w.Expr, next)

	w.Expr = ref

	return append(toCommands(cmds), w)
}
