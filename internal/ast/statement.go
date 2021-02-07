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
	String() string
	lValue()
}

type LTuple struct {
	Args []LValue
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

func (l *LetStatement) command()   {}
func (l *LetStatement) statement() {}

func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", l.LValue.String(), l.Expr.String())
}

type ReturnStatement struct {
	Expr Expression
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", r.Expr.String())
}
func (r *ReturnStatement) command()   {}
func (r *ReturnStatement) statement() {}
