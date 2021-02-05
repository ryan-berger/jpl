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
	Type     string
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

	return fmt.Sprintf("{%s}", strings.Join(strs, ","))
}
func (b TupleBinding) binding() {}
