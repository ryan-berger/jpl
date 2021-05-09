package optimizer

import (
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

func searchStmt(d *defUse, statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		arg := stmt.LValue.(*ast.Variable)
		d.recordDef(arg.Variable)
		searchUses(d, stmt, stmt.Expr)
	case *ast.AssertStatement:
		searchUses(d, stmt, stmt.Expr)
	case *ast.ReturnStatement:
		searchUses(d, stmt, stmt.Expr)
	}
}

func searchCmd(d *defUse, command ast.Command) {
	switch cmd := command.(type) {
	case ast.Statement:
		searchStmt(d, cmd)
	case *ast.Show:
		searchUses(d, cmd, cmd.Expr)
	case *ast.Time:
		searchCmd(d, cmd)
	case *ast.Function:
		def := makeDefUse(d)
		d.children[cmd] = def
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

func shouldRemove(n ast.Node, use *defUse) bool {
	if let, ok := n.(*ast.LetStatement); ok {
		variable := let.LValue.(*ast.Variable).Variable
		uses := use.getUses(variable)
		if len(uses) == 0 {
			use.clearUse(let)
			return true
		}
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
