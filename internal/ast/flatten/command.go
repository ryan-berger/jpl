package flatten

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

func flattenCommand(command ast.Command, next nexter) []ast.Command {
	switch cmd := command.(type) {
	case *ast.Read:
		return []ast.Command{cmd}
	case *ast.Print:
		return []ast.Command{cmd}
	case *ast.Write:
		return flattenWrite(cmd, next)
	case *ast.Time:
		return flattenTime(cmd, next)
	case *ast.Show:
		return flattenShow(cmd, next)
	case *ast.Function:
		return []ast.Command{flattenFunction(cmd, next)}
	case ast.Statement:
		return toCommands(flattenStatement(cmd, next))
	default:
		panic("oops, type not supported")
	}
	return nil
}

func flattenFunction(fn *ast.Function, next nexter) ast.Command {
	stmts := fn.Statements
	var flattened []ast.Statement
	for _, s := range stmts {
		flattened = append(flattened, flattenStatement(s, next)...)

		if size := len(flattened); size != 0 {
			if isReturn(flattened[size-1]) { // exit early, we've hit our return
				fn.Statements = flattened
				return fn
			}
		}
	}

	if size := len(flattened); size != 0 && isReturn(flattened[size-1]) { // add a return at last since there is none
		name := next()
		l := dsl.Let(
			dsl.LIdent(name), dsl.Tuple())
		ret := dsl.Return(dsl.Ident(name))

		flattened = append(flattened, l, ret)
	}

	fn.Statements = flattened
	return fn
}

func flattenShow(sh *ast.Show, next nexter) []ast.Command {
	ref, cmds := expansionAndLet(sh.Expr, next)

	sh.Expr = ref

	var combined []ast.Command
	combined = append(combined, toCommands(cmds)...)
	combined = append(combined, sh)
	return combined
}

func flattenTime(time *ast.Time, next nexter) []ast.Command {
	startRef := next()
	start := dsl.Let(
		dsl.LIdent(startRef),
		dsl.Call("get_time"))

	cmds := flattenCommand(time.Command, next)

	size := len(cmds)
	if size != 0 {
		if _, ok := cmds[size-1].(*ast.ReturnStatement); ok {
			return append([]ast.Command{start}, cmds...)
		}
	}

	endRef := next()
	end := dsl.Let(dsl.LIdent(endRef),
		dsl.Call("get_time"))

	subRef := next()
	sub := dsl.Let(
		dsl.LIdent(subRef),
		dsl.Infix("-",
			dsl.Ident(endRef),
			dsl.Ident(startRef)))

	p := dsl.Print("time: ")
	s := dsl.Show(dsl.Ident(subRef))

	commands := []ast.Command{start}
	commands = append(commands, cmds...)
	commands = append(commands, end, sub, p, s)
	return commands
}

func flattenWrite(w *ast.Write, next nexter) []ast.Command {
	ref, cmds := expansionAndLet(w.Expr, next)

	w.Expr = ref

	return append(toCommands(cmds), w)
}
