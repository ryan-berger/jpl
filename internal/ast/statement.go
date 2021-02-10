package ast

import (
	"fmt"
	"strings"
)

type Statement interface {
	Command
	statement()
}

type LValue interface {
	SExpr
	String() string
	lValue()
}

type LTuple struct {
	Args []LValue
}

func (l *LTuple) SExpr() string {
	panic("implement me")
}

func (l *LTuple) String() string {
	lVals := make([]string, len(l.Args))
	for i := 0; i < len(l.Args); i++ {
		lVals[i] = l.Args[i].String()
	}
	return fmt.Sprintf("{%s}", strings.Join(lVals, ", "))
}
func (l *LTuple) lValue() {}

type LetStatement struct {
	LValue LValue
	Expr   Expression
}

func (l *LetStatement) SExpr() string {
	return fmt.Sprintf("(LetStmt %s %s)", l.LValue.String(), l.Expr.SExpr())
}

func (l *LetStatement) command()   {}
func (l *LetStatement) statement() {}

func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", l.LValue, l.Expr)
}

type ReturnStatement struct {
	Expr Expression
}

func (r *ReturnStatement) SExpr() string {
	return fmt.Sprintf("(ReturnStmt %s)", r.Expr.SExpr())
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", r.Expr.String())
}
func (r *ReturnStatement) command()   {}
func (r *ReturnStatement) statement() {}

type AssertStatement struct {
	Expr    Expression
	Message string
}

func (a *AssertStatement) SExpr() string {
	return fmt.Sprintf("(AssertStmt %s %s)", a.Expr.SExpr(), a.Message)
}

func (a *AssertStatement) String() string {
	return fmt.Sprintf("assert %s , %s", a.Expr.String(), a.Message)
}

func (a *AssertStatement) command()   {}
func (a *AssertStatement) statement() {}
