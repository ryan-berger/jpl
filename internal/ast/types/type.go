package types

import (
	"fmt"
	"strings"
)

type Type interface {
	Equal(other Type) bool
	Size() int
	String() string
	SExpr() string
}

type boolean struct{}

func (b *boolean) Size() int {
	return 4
}

func (b *boolean) String() string {
	return "bool"
}

func (b *boolean) Equal(other Type) bool {
	_, ok := other.(*boolean)
	return ok
}

func (b *boolean) SExpr() string {
	return "BoolType"
}

type integer struct{}

func (i *integer) Size() int {
	return 8
}

func (i *integer) Equal(other Type) bool {
	_, ok := other.(*integer)
	return ok
}

func (i *integer) String() string {
	return "int"
}

func (i *integer) SExpr() string {
	return "IntType"
}

type float struct{}

func (f *float) Size() int {
	return 8
}

func (f *float) Equal(other Type) bool {
	_, ok := other.(*float)
	return ok
}

func (f *float) String() string {
	return "float"
}

func (f *float) SExpr() string {
	return "FloatType"
}

type Array struct {
	Inner Type
	Rank  int
}

func (a *Array) Size() int {
	return (a.Rank + 1) * 8
}

func (a *Array) Equal(other Type) bool {
	arr, ok := other.(*Array)
	return ok && a.Rank == arr.Rank && a.Inner.Equal(arr.Inner)
}

func (a *Array) String() string {
	if a == Pict {
		return "pict" // special case
	}

	b := make([]byte, a.Rank-1)
	for i := 0; i < len(b); i++ {
		b[i] = ','
	}
	return fmt.Sprintf("%s[%s]", a.Inner, string(b))
}

func (a *Array) SExpr() string {
	return fmt.Sprintf("(ArrayType %s rank=%d)",
		a.Inner.SExpr(), a.Rank)
}

type Tuple struct {
	Types []Type
}

func (t *Tuple) Size() int {
	return 0
}

func (t *Tuple) Equal(other Type) bool {
	tup, ok := other.(*Tuple)
	if !ok {
		return false
	}

	if len(t.Types) != len(tup.Types) {
		return false
	}

	for i, typ := range t.Types {
		if !typ.Equal(tup.Types[i]) {
			return false
		}
	}

	return true
}

func (t *Tuple) String() string {
	strs := make([]string, len(t.Types))
	for i, typ := range t.Types {
		strs[i] = typ.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (t *Tuple) SExpr() string { return "" }

type str struct{}

func (s *str) Equal(other Type) bool {
	_, ok := other.(*str)
	return ok
}

func (s *str) Size() int {
	//TODO implement me
	panic("implement me")
}

func (s *str) String() string {
	return "str"
}

func (s *str) SExpr() string {
	return "StrType"
}
