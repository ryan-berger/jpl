package ast

import (
	"fmt"
	"strings"
)

type Expression interface {
	String() string
	expression()
}

// IntExpression
type IntExpression struct {
	Val int64
}

func (i *IntExpression) String() string {
	return fmt.Sprintf("%d", i.Val)
}
func (i *IntExpression) expression() {}

// IdentifierExpression
type IdentifierExpression struct {
	Identifier string
}

func (i *IdentifierExpression) String() string {
	return i.Identifier
}
func (i *IdentifierExpression) expression() {}

type CallExpression struct {
	Identifier string
	Arguments  []Expression
}

func (c *CallExpression) String() string {
	strs := make([]string, len(c.Arguments))
	for i, expr := range c.Arguments {
		strs[i] = expr.String()
	}

	return fmt.Sprintf("%s(%s)", c.Identifier, strings.Join(strs, ", "))
}
func (c *CallExpression) expression() {}

// FloatExpression
type FloatExpression struct {
	Val float64
}

func (f *FloatExpression) String() string {
	return fmt.Sprintf("%f", f.Val)
}
func (f *FloatExpression) expression() {}

type BooleanExpression struct {
	Val bool
}

func (b *BooleanExpression) String() string {
	if b.Val {
		return "true"
	}
	return "false"
}
func (b *BooleanExpression) expression() {}

type TupleExpression struct {
	Expressions []Expression
}

func (t *TupleExpression) String() string {
	strs := make([]string, len(t.Expressions))
	for i, expr := range t.Expressions {
		strs[i] = expr.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (t *TupleExpression) expression() {}

type IfExpression struct {
	Condition   Expression
	Consequence Expression
	Otherwise   Expression
}

func (i *IfExpression) String() string {
	return fmt.Sprintf("if %s then %s else %s",
		i.Condition.String(), i.Consequence.String(), i.Otherwise.String())
}
func (i *IfExpression) expression() {}

type OpBinding struct {
	Variable string
	Expr     Expression
}

func (o *OpBinding) String() string {
	return fmt.Sprintf("%s : %s", o.Variable, o.Expr)
}

type ArrayTransform struct {
	OpBindings []OpBinding
	Expr       Expression
}

func (a *ArrayTransform) String() string {
	bindings := make([]string, len(a.OpBindings))
	for i, b := range a.OpBindings {
		bindings[i] = b.String()
	}

	return fmt.Sprintf("array[%s] %s", strings.Join(bindings, ", "), a.Expr)
}

func (a *ArrayTransform) expression() {}

type SumTransform struct {
	OpBindings []OpBinding
	Expr       Expression
}

func (s *SumTransform) String() string {
	bindings := make([]string, len(s.OpBindings))
	for i, b := range s.OpBindings {
		bindings[i] = b.String()
	}
	return fmt.Sprintf("array[%s] %s", strings.Join(bindings, ", "), s.Expr)
}
func (s *SumTransform) expression() {}

type InfixExpression struct {
	Left  Expression
	Right Expression
	Op    string
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Op, i.Right)
}
func (i *InfixExpression) expression() {}

type PrefixExpression struct {
	Op   string
	Expr Expression
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Expr)
}

func (p *PrefixExpression) expression() {}
