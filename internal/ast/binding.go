package ast

import (
	"fmt"
	"strings"

	"github.com/ryan-berger/jpl/internal/types"
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
	panic("implement me")
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



type BasicType int

func (b BasicType) String() string {
	if b == Int {
		return "int"
	}
	if b == Boolean {
		return "bool"
	}
	return "float"
}

func (b BasicType) typ() {}

const (
	Int BasicType = iota
	Float
	Boolean
)

type ArrType struct {
	Type types.Type
	Rank int
}

func (a *ArrType) String() string {
	commas := make([]byte, a.Rank)
	for i := 0; i < a.Rank-1; i++ {
		commas[i] = ','
	}
	return fmt.Sprintf("%s[%s]", a.Type, string(commas))
}
func (a *ArrType) typ() {}

type TupleType struct {
	Types []types.Type
}

func (t *TupleType) String() string {
	strs := make([]string, len(t.Types))
	for i, t := range t.Types {
		strs[i] = t.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}
func (t *TupleType) typ() {}
