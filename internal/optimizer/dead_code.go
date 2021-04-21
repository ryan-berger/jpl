package optimizer

import "github.com/ryan-berger/jpl/internal/ast"

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

func DeadCode(p ast.Program) ast.Program {
	p = deadStmts(p)
	return p
}
