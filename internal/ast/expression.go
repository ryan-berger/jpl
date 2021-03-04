package ast

import (
	"fmt"
	"strings"
)

type Expression interface {
	Node
	expression()
}

// IntExpression
type IntExpression struct {
	Val  int64
	Type Type
	Location
}

func (i *IntExpression) SExpr() string {
	return fmt.Sprintf("(IntExpr %d)", i.Val)
}

func (i *IntExpression) String() string {
	return fmt.Sprintf("%d", i.Val)
}
func (i *IntExpression) expression() {}

// IdentifierExpression
type IdentifierExpression struct {
	Identifier string
	Location
}

func (i *IdentifierExpression) SExpr() string {
	return fmt.Sprintf("(VarExpr %s)", i.Identifier)
}

func (i *IdentifierExpression) String() string {
	return i.Identifier
}
func (i *IdentifierExpression) expression() {}

type CallExpression struct {
	Identifier string
	Arguments  []Expression
	Location
}

func (c *CallExpression) SExpr() string {
	strs := make([]string, len(c.Arguments))
	for i, expr := range c.Arguments {
		strs[i] = expr.SExpr()
	}

	return fmt.Sprintf("(CallExpr %s %s)", c.Identifier, strings.Join(strs, " "))
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
	Location
}

func (f *FloatExpression) SExpr() string {
	return fmt.Sprintf("(FloatExpr %d)", int64(f.Val))
}

func (f *FloatExpression) String() string {
	return fmt.Sprintf("%f", f.Val)
}
func (f *FloatExpression) expression() {}

type BooleanExpression struct {
	Val bool
	Location
}

func (b *BooleanExpression) SExpr() string {
	panic("implement me")
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
	Location
}

func (t *TupleExpression) SExpr() string {
	panic("implement me")
}

func (t *TupleExpression) String() string {
	strs := make([]string, len(t.Expressions))
	for i, expr := range t.Expressions {
		strs[i] = expr.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (t *TupleExpression) expression() {}

type ArrayExpression struct {
	Expressions []Expression
	Location
}

func (a *ArrayExpression) SExpr() string {
	return ""
}

func (a *ArrayExpression) String() string {
	strs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		strs[i] = e.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func (a *ArrayExpression) expression() {}

type ArrayRefExpression struct {
	Array   Expression
	Indexes []Expression
	Location
}

func (a *ArrayRefExpression) SExpr() string {
	return ""
}
func (a *ArrayRefExpression) String() string {
	strs := make([]string, len(a.Indexes))
	for i, idx := range a.Indexes {
		strs[i] = fmt.Sprintf("%s", idx)
	}
	return fmt.Sprintf("%s[%s]", a.Array, strings.Join(strs, ", "))
}
func (a *ArrayRefExpression) expression() {}

type TupleRefExpression struct {
	Tuple Expression
	Index Expression
	Location
}

func (t *TupleRefExpression) SExpr() string {
	panic("implement me")
}
func (t *TupleRefExpression) String() string {
	return fmt.Sprintf("%s{%s}", t.Tuple, t.Index)
}
func (t *TupleRefExpression) expression() {}

type IfExpression struct {
	Condition   Expression
	Consequence Expression
	Otherwise   Expression
	Location
}

func (i *IfExpression) SExpr() string {
	panic("implement me")
}

func (i *IfExpression) String() string {
	return fmt.Sprintf("if %s then %s else %s",
		i.Condition.String(), i.Consequence.String(), i.Otherwise.String())
}
func (i *IfExpression) expression() {}

type OpBinding struct {
	Variable string
	Expr     Expression
	Location
}

func (o *OpBinding) String() string {
	return fmt.Sprintf("%s : %s", o.Variable, o.Expr)
}

type ArrayTransform struct {
	OpBindings []OpBinding
	Expr       Expression
	Location
}

func (a *ArrayTransform) SExpr() string {
	panic("implement me")
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
	Location
}

func (s *SumTransform) SExpr() string {
	panic("implement me")
}

func (s *SumTransform) String() string {
	bindings := make([]string, len(s.OpBindings))
	for i, b := range s.OpBindings {
		bindings[i] = b.String()
	}
	return fmt.Sprintf("sum[%s] %s", strings.Join(bindings, ", "), s.Expr)
}
func (s *SumTransform) expression() {}

type InfixExpression struct {
	Left  Expression
	Right Expression
	Op    string
	Location
}

func (i *InfixExpression) SExpr() string {
	panic("shouldn't run")
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Op, i.Right)
}
func (i *InfixExpression) expression() {}

type PrefixExpression struct {
	Op   string
	Expr Expression
	Location
}

func (p *PrefixExpression) SExpr() string {
	return p.Expr.SExpr()
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Expr)
}

func (p *PrefixExpression) expression() {}
