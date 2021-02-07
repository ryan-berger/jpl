package ast

import (
	"fmt"
	"strings"
)

type Binding interface {
	String() string
	binding()
}

type TypeBind struct {
	Argument Argument
	Type     Type
}

func (b *TypeBind) String() string {
	return fmt.Sprintf("%s : %s", b.Argument.String(), b.Type)
}
func (b *TypeBind) binding() {}

type TupleBinding []Binding

func (b TupleBinding) String() string {
	strs := make([]string, len(b))
	for i, b := range b {
		strs[i] = b.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}
func (b TupleBinding) binding() {}

type Type interface {
	String() string
	typ()
}

type BasicType int

func (b BasicType) String() string {
	if b == Int {
		return "int"
	}
	return "float"
}

func (b BasicType) typ() {}

const (
	Int BasicType = iota
	Float
)

type ArrType struct {
	Type Type
	Rank int
}

func (a *ArrType) String() string {
	commas := make([]byte, a.Rank)
	for i := 0; i < a.Rank - 1; i++ {
		commas[i] = ','
	}
	return fmt.Sprintf("%s[%s]", a.Type.String(), string(commas))
}
func (a *ArrType) typ() {}

type TupleType struct {
	Types []Type
}

func (t *TupleType) String() string {
	strs := make([]string, len(t.Types))
	for i, t := range t.Types {
		strs[i] = t.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (t *TupleType) typ() {
	panic("implement me")
}

