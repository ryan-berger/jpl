package ast

import (
	"fmt"
	"strings"

	"github.com/ryan-berger/jpl/internal/types"
)

type Expression interface {
	Node
	Typ() types.Type
	expression()
}

// IntExpression
type IntExpression struct {
	Val  int64
	Type types.Type
	Location
}

func (i *IntExpression) SExpr() string {
	if i.Type != nil {
		return fmt.Sprintf("(IntExpr %s %d)", i.Type.SExpr(), i.Val)
	}
	return fmt.Sprintf("(IntExpr %d)", i.Val)
}

func (i *IntExpression) String() string {
	return fmt.Sprintf("%d", i.Val)
}

func (i *IntExpression) Typ() types.Type {
	return i.Type
}

func (i *IntExpression) expression() {}

// IdentifierExpression
type IdentifierExpression struct {
	Identifier string
	Type       types.Type
	Location
}

func (i *IdentifierExpression) SExpr() string {
	if i.Type != nil {
		return fmt.Sprintf("(VarExpr %s %s)", i.Type.SExpr(), i.Identifier)
	}
	return fmt.Sprintf("(VarExpr %s)", i.Identifier)
}

func (i *IdentifierExpression) String() string {
	return i.Identifier
}

func (i *IdentifierExpression) Typ() types.Type {
	return i.Type
}

func (i *IdentifierExpression) expression() {}

type CallExpression struct {
	Identifier string
	Arguments  []Expression
	Type       types.Type
	Location
}

func (c *CallExpression) SExpr() string {
	strs := make([]string, len(c.Arguments))
	for i, expr := range c.Arguments {
		strs[i] = expr.SExpr()
	}

	if c.Type != nil {
		return fmt.Sprintf("(CallExpr %s %s %s)", c.Type.SExpr(), c.Identifier, strings.Join(strs, " "))
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

func (c *CallExpression) Typ() types.Type {
	return c.Type
}

func (c *CallExpression) expression() {}

// FloatExpression
type FloatExpression struct {
	Val  float64
	Type types.Type
	Location
}

func (f *FloatExpression) SExpr() string {
	if f.Type != nil {
		return fmt.Sprintf("(FloatExpr %s %d)", f.Type.SExpr(), int64(f.Val))
	}
	return fmt.Sprintf("(FloatExpr %d)", int64(f.Val))
}

func (f *FloatExpression) String() string {
	return fmt.Sprintf("%f", f.Val)
}

func (f *FloatExpression) Typ() types.Type {
	return f.Type
}

func (f *FloatExpression) expression() {}

type BooleanExpression struct {
	Val  bool
	Type types.Type
	Location
}

func (b *BooleanExpression) SExpr() string {
	if b.Type != nil {
		return fmt.Sprintf("(VarExpr %s %t)", b.Type.SExpr(), b.Val)
	}
	return fmt.Sprintf("(VarExpr %t)", b.Val)
}

func (b *BooleanExpression) String() string {
	if b.Val {
		return "true"
	}
	return "false"
}

func (b *BooleanExpression) Typ() types.Type {
	return b.Type
}

func (b *BooleanExpression) expression() {}

type TupleExpression struct {
	Expressions []Expression
	Type        types.Type
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

func (t *TupleExpression) Typ() types.Type {
	return t.Type
}

func (t *TupleExpression) expression() {}

type ArrayExpression struct {
	Expressions []Expression
	Type        types.Type
	Location
}

func (a *ArrayExpression) SExpr() string {
	strs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		strs[i] = e.SExpr()
	}
	if a.Type != nil {
		return fmt.Sprintf("(ArrConExpr %s %s)",
			a.Type.SExpr(), strings.Join(strs, " "))
	}
	return fmt.Sprintf("(ArrConExpr %s)", strings.Join(strs, " "))
}

func (a *ArrayExpression) String() string {
	strs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		strs[i] = e.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func (a *ArrayExpression) Typ() types.Type {
	return a.Type
}

func (a *ArrayExpression) expression() {}

type ArrayRefExpression struct {
	Array   Expression
	Indexes []Expression
	Type    types.Type
	Location
}

func (a *ArrayRefExpression) SExpr() string {
	indexStrs := make([]string, len(a.Indexes))
	for i, a := range a.Indexes {
		indexStrs[i] = a.SExpr()
	}

	if a.Type != nil {
		return fmt.Sprintf("(ArrIndexExpr %s %s %s)",
			a.Type.SExpr(), a.Array.SExpr(),
			strings.Join(indexStrs, " "))
	}
	return fmt.Sprintf("")
}
func (a *ArrayRefExpression) String() string {
	strs := make([]string, len(a.Indexes))
	for i, idx := range a.Indexes {
		strs[i] = fmt.Sprintf("%s", idx)
	}
	return fmt.Sprintf("%s[%s]", a.Array, strings.Join(strs, ", "))
}

func (a *ArrayRefExpression) Typ() types.Type {
	return a.Type
}

func (a *ArrayRefExpression) expression() {}

type TupleRefExpression struct {
	Tuple Expression
	Index int64
	Type  types.Type
	Location
}

func (t *TupleRefExpression) SExpr() string {
	panic("implement me")
}

func (t *TupleRefExpression) String() string {
	return fmt.Sprintf("%s{%d}", t.Tuple, t.Index)
}

func (t *TupleRefExpression) Typ() types.Type {
	return t.Type
}

func (t *TupleRefExpression) expression() {}

type IfExpression struct {
	Condition   Expression
	Consequence Expression
	Otherwise   Expression
	Type        types.Type
	Location
}

func (i *IfExpression) SExpr() string {
	if i.Type != nil {
		return fmt.Sprintf("(IteExpr %s %s %s %s)",
			i.Type.SExpr(),
			i.Condition.SExpr(),
			i.Consequence.SExpr(),
			i.Otherwise.SExpr())
	}
	return fmt.Sprintf("(IteExpr %s %s %s)",
		i.Condition.SExpr(),
		i.Consequence.SExpr(),
		i.Otherwise.SExpr())
}

func (i *IfExpression) String() string {
	return fmt.Sprintf("if %s then %s else %s",
		i.Condition.String(), i.Consequence.String(), i.Otherwise.String())
}

func (i *IfExpression) Typ() types.Type {
	return i.Type
}

func (i *IfExpression) expression() {}

type OpBinding struct {
	Variable string
	Expr     Expression
	Type     types.Type
	Location
}

func (o *OpBinding) SExpr() string {
	return fmt.Sprintf("(%s %s)", o.Variable, o.Expr.SExpr())
}

func (o *OpBinding) String() string {
	return fmt.Sprintf("%s : %s", o.Variable, o.Expr)
}

type ArrayTransform struct {
	OpBindings []OpBinding
	Expr       Expression
	Type       types.Type
	Location
}

func (a *ArrayTransform) SExpr() string {
	bindingStrs := make([]string, len(a.OpBindings))
	for i, b := range a.OpBindings {
		bindingStrs[i] = b.SExpr()
	}

	if a.Type != nil {
		return fmt.Sprintf("(ArrayExpr %s %s %s)",
			a.Type.SExpr(),
			strings.Join(bindingStrs, " "),
			a.Expr.SExpr())
	}

	return fmt.Sprintf("(ArrayExpr %s)", a.Expr.SExpr())
}

func (a *ArrayTransform) Typ() types.Type {
	return a.Type
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
	Type       types.Type
	Location
}

func (s *SumTransform) SExpr() string {
	bindingStrs := make([]string, len(s.OpBindings))
	for i, b := range s.OpBindings {
		bindingStrs[i] = b.SExpr()
	}

	if s.Type != nil {
		return fmt.Sprintf("(SumExpr %s %s %s)",
			s.Type.SExpr(),
			strings.Join(bindingStrs, " "),
			s.Expr.SExpr())
	}

	return fmt.Sprintf("(SumExpr %s)", s.Expr.SExpr())
}

func (s *SumTransform) String() string {
	bindings := make([]string, len(s.OpBindings))
	for i, b := range s.OpBindings {
		bindings[i] = b.String()
	}
	return fmt.Sprintf("sum[%s] %s", strings.Join(bindings, ", "), s.Expr)
}

func (s *SumTransform) Typ() types.Type {
	return s.Type
}
func (s *SumTransform) expression() {}

type InfixExpression struct {
	Left  Expression
	Right Expression
	Op    string
	Type  types.Type
	Location
}

func (i *InfixExpression) SExpr() string {
	if i.Type != nil {
		return fmt.Sprintf("(BinopExpr %s %s %s %s)",
			i.Type.SExpr(),
			i.Left.SExpr(),
			i.Op,
			i.Right.SExpr())
	}
	return fmt.Sprintf("(BinopExpr %s %s %s)",
		i.Left.SExpr(),
		i.Op,
		i.Right.SExpr())
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Op, i.Right)
}

func (i *InfixExpression) Typ() types.Type {
	return i.Type
}

func (i *InfixExpression) expression() {}

type PrefixExpression struct {
	Op   string
	Expr Expression
	Type types.Type
	Location
}

func (p *PrefixExpression) Typ() types.Type {
	return p.Type
}

func (p *PrefixExpression) SExpr() string {
	if p.Type != nil {
		return fmt.Sprintf("(UnopExpr %s %s %s)",
			p.Type.SExpr(), p.Op, p.Expr.SExpr())
	}
	return fmt.Sprintf("(UnopExpr %s %s)",
		p.Op, p.Expr.SExpr())
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Expr)
}

func (p *PrefixExpression) expression() {}
