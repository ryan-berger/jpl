package ast

import (
	"fmt"
	"strings"

	types "github.com/ryan-berger/jpl/internal/ast/types"
)

type Binding interface {
	Node
	binding()
}

type TypeBind struct {
	Argument Argument
	Type     types.Type
	Location
}

func (b *TypeBind) SExpr() string {
	return fmt.Sprintf("(Binding %s %s)",
		b.Argument.SExpr(),
		b.Type.SExpr())
}

func (b *TypeBind) String() string {
	return fmt.Sprintf("%s : %s", b.Argument, b.Type)
}
func (b *TypeBind) binding() {}

type TupleBinding struct {
	Bindings []Binding
	Location
}

func (b *TupleBinding) SExpr() string {
	panic("implement me")
}

func (b *TupleBinding) String() string {
	strs := make([]string, len(b.Bindings))
	for i, b := range b.Bindings {
		strs[i] = b.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}
func (b *TupleBinding) binding() {}
