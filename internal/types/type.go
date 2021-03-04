package types

import (
	"fmt"
	"strings"
)

type Type interface {
	Equal(other Type) bool
	String() string
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

type integer struct{}

func (i *integer) Equal(other Type) bool {
	_, ok := other.(*integer)
	return ok
}

func (i *integer) String() string {
	return "int"
}

type float struct{}

func (f *float) Equal(other Type) bool {
	_, ok := other.(*float)
	return ok
}

func (f *float) String() string {
	return "float"
}

type Array struct {
	Inner Type
	Rank  int
}

func (a *Array) Equal(other Type) bool {
	arr, ok := other.(*Array)
	return ok && a.Rank == arr.Rank && a.Inner.Equal(arr.Inner)
}

func (a *Array) String() string {
	if a == Pict {
		return "pict" // special case
	}

	b := make([]byte, a.Rank - 1)
	for i := 0; i < len(b); i++ {
		b[i] = ','
	}
	return fmt.Sprintf("%s[%s]", a.Inner, string(b))
}

type Tuple struct {
	Types []Type
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