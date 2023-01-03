package optimizer

import (
	"fmt"
	"sort"

	"github.com/ryan-berger/jpl/internal/ast"
)

type deadCode int

const (
	nop deadCode = iota
	skip
	fullStop
)

func deadCodeExpr(expression ast.Expression) ast.Expression {
	switch expr := expression.(type) {
	case *ast.CallExpression:
		for i := 0; i < len(expr.Arguments); i++ {
			expr.Arguments[i] = deadCodeExpr(expr.Arguments[i])
		}
	case *ast.SumTransform:
		for _, bin := range expr.OpBindings {
			bin.Expr = deadCodeExpr(bin.Expr)
		}
		expr.Expr = deadCodeExpr(expr.Expr)
	case *ast.ArrayTransform:
		for _, bin := range expr.OpBindings {
			bin.Expr = deadCodeExpr(bin.Expr)
		}
		expr.Expr = deadCodeExpr(expr.Expr)
	case *ast.InfixExpression:
		expr.Left = deadCodeExpr(expr.Left)
		expr.Right = deadCodeExpr(expr.Right)
	case *ast.PrefixExpression:
		expr.Expr = deadCodeExpr(expr.Expr)
	case *ast.IfExpression:
		if v, ok := expr.Condition.(*ast.BooleanExpression); ok {
			if v.Val {
				return expr.Consequence
			}
			return expr.Otherwise
		}
	}
	return expression
}

func deadCodeStmt(statement ast.Statement) deadCode {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		stmt.Expr = deadCodeExpr(stmt.Expr)
	case *ast.ReturnStatement:
		stmt.Expr = deadCodeExpr(stmt.Expr)
	case *ast.AssertStatement:
		v, ok := stmt.Expr.(*ast.BooleanExpression)
		if ok {
			if v.Val {
				return skip
			}
			return fullStop
		}
		stmt.Expr = deadCodeExpr(stmt.Expr)
	}
	return nop
}

func deadCodeCmd(command ast.Command) deadCode {
	switch cmd := command.(type) {
	case ast.Statement:
		return deadCodeStmt(cmd)
	case *ast.Time:
		return deadCodeCmd(cmd.Command)
	case *ast.Write:
		cmd.Expr = deadCodeExpr(cmd.Expr)
	case *ast.Show:
		cmd.Expr = deadCodeExpr(cmd.Expr)
	case *ast.Function:
		var stmts []ast.Statement
		for _, stmt := range cmd.Statements {
			code := deadCodeStmt(stmt)
			switch code {
			case skip:
				continue
			case fullStop:
				last := cmd.Statements[len(cmd.Statements)-1]
				stmts = append(stmts, stmt, last)
				goto end
			}

			stmts = append(stmts, stmt)
		}
	end:
		cmd.Statements = stmts
	}

	return nop
}

func deadStmts(p ast.Program) ast.Program {
	var cmds []ast.Command

	for _, cmd := range p {
		code := deadCodeCmd(cmd)
		switch code {
		case skip:
			continue
		case fullStop:
			cmds = append(cmds, cmd)
			if ret, isReturn := p[len(p)-1].(*ast.ReturnStatement); isReturn {
				cmds = append(cmds, ret)
			}
			goto end
		}
		cmds = append(cmds, cmd)
	}
end:
	return cmds
}

func searchUses(du *defUse, parent ast.Node, expr ast.Expression) {
	switch exp := expr.(type) {
	case *ast.IdentifierExpression:
		du.recordUse(exp.Identifier, parent)
	case *ast.ArrayExpression:
		for _, e := range exp.Expressions {
			searchUses(du, parent, e)
		}
	case *ast.TupleExpression:
		for _, e := range exp.Expressions {
			searchUses(du, parent, e)
		}
	case *ast.TupleRefExpression:
		searchUses(du, parent, exp.Tuple)
	case *ast.CallExpression:
		for _, arg := range exp.Arguments {
			searchUses(du, parent, arg)
		}
	case *ast.ArrayTransform:
		def := makeDefUse(du)
		for _, arg := range exp.OpBindings {
			def.recordDef(arg.Variable)
			searchUses(def, parent, arg.Expr)
		}
		searchUses(def, parent, exp.Expr)
	case *ast.SumTransform:
		def := makeDefUse(du)
		for _, arg := range exp.OpBindings {
			def.recordDef(arg.Variable)
			searchUses(def, parent, arg.Expr)
		}
		searchUses(def, parent, exp.Expr)
	case *ast.IfExpression:
		searchUses(du, parent, exp.Condition)
		searchUses(du, parent, exp.Otherwise)
		searchUses(du, parent, exp.Consequence)
	case *ast.InfixExpression:
		searchUses(du, parent, exp.Left)
		searchUses(du, parent, exp.Right)
	case *ast.PrefixExpression:
		searchUses(du, parent, exp.Expr)
	}
}

func recordDefs(lval ast.LValue, d *defUse) {
	switch v := lval.(type) {
	case *ast.Variable:
		d.recordDef(v.Variable)
	case *ast.VariableArr:
		d.recordDef(v.Variable)
		for _, arr := range v.Variables {
			d.recordDef(arr)
		}
	case *ast.LTuple:
		for _, t := range v.Args {
			recordDefs(t, d)
		}
	}
}

func searchStmt(d *defUse, statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		recordDefs(stmt.LValue, d)
		searchUses(d, stmt, stmt.Expr)
	case *ast.AssertStatement:
		searchUses(d, stmt, stmt.Expr)
	case *ast.ReturnStatement:
		searchUses(d, stmt, stmt.Expr)
	}
}

func searchBinding(d *defUse, binding ast.Binding) {
	switch bind := binding.(type) {
	case *ast.TypeBind:
		switch arg := bind.Argument.(type) {
		case *ast.Variable:
			fmt.Println(arg.Variable)
			d.recordDef(arg.Variable)
		case *ast.VariableArr:
			d.recordDef(arg.Variable)
			for _, v := range arg.Variables {
				d.recordDef(v)
			}
		}
	case *ast.TupleBinding:
		for _, b := range bind.Bindings {
			searchBinding(d, b)
		}
	}
}

func searchCmd(d *defUse, command ast.Command) {
	switch cmd := command.(type) {
	case ast.Statement:
		searchStmt(d, cmd)
	case *ast.Show:
		searchUses(d, cmd, cmd.Expr)
	case *ast.Time:
		searchCmd(d, cmd.Command)
	case *ast.Function:
		def := makeDefUse(d)
		d.children[cmd] = def

		for _, bind := range cmd.Bindings {
			searchBinding(def, bind)
		}

		fmt.Println(d.graph)

		for _, stmt := range cmd.Statements {
			searchStmt(def, stmt)
		}
	}
}

func makeDefUse(parent *defUse) *defUse {
	return &defUse{
		graph:        make(map[string]map[ast.Node]bool),
		reverseGraph: make(map[ast.Node]map[string]bool),
		parent:       parent,
		children:     make(map[*ast.Function]*defUse),
	}
}

func buildDefUse(p ast.Program) *defUse {
	d := makeDefUse(nil)

	for _, cmd := range p {
		searchCmd(d, cmd)
	}

	return d
}

func hasNoUses(l ast.LValue, use *defUse) bool {
	switch v := l.(type) {
	case *ast.Variable:
		return len(use.getUses(v.Variable)) == 0
	case *ast.VariableArr:
		canClear := len(use.getUses(v.Variable)) == 0
		for i := 0; i < len(v.Variables) && canClear; i++ {
			canClear = canClear && len(use.getUses(v.Variables[i])) == 0
		}
		return canClear
	case *ast.LTuple:
		canClear := true
		for i := 0; i < len(v.Args) && canClear; i++ {
			fmt.Printf("%s has uses? %v", v.Args[i], hasNoUses(v.Args[i], use))
			canClear = canClear && hasNoUses(v.Args[i], use)
		}
		return canClear
	default:
		panic("unreachable")
	}
}

func shouldRemove(n ast.Node, use *defUse) bool {
	switch stmt := n.(type) {
	case *ast.LetStatement:
		if hasNoUses(stmt.LValue, use) {
			fmt.Printf("%s has no uses\n", stmt.LValue.String())
			use.clearUse(stmt)
			return true
		}
	case *ast.AssertStatement:
		val, ok := stmt.Expr.(*ast.BooleanExpression)
		use.clearUse(stmt)
		return ok && val.Val
	}
	return false
}

func removeUnused(p ast.Program, use *defUse) ast.Program {
	var cmds ast.Program
	for i := len(p) - 1; i >= 0; i-- {
		if shouldRemove(p[i], use) {
			continue
		}

		if fn, ok := p[i].(*ast.Function); ok {
			var fnStmts []ast.Statement
			fnUse := use.children[fn]
			for j := len(fn.Statements) - 1; j >= 0; j-- {
				stmt := fn.Statements[j]
				if shouldRemove(stmt, fnUse) {
					continue
				}
				fnStmts = append(fnStmts, stmt)
			}
			sort.Slice(fnStmts, func(i, j int) bool {
				return j < i
			})
			fn.Statements = fnStmts
		}
		cmds = append(cmds, p[i])
	}

	// reverse
	sort.Slice(cmds, func(i, j int) bool {
		return j < i
	})

	return cmds
}

func DeadCode(p ast.Program) ast.Program {
	p = deadStmts(p)
	du := buildDefUse(p)
	return removeUnused(p, du)
}
