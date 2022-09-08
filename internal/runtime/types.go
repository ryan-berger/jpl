package runtime

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast/types"
)

var String = &str{}

type str struct{}

func (s *str) Equal(other types.Type) bool {
	return other == String
}
func (s *str) Size() int {
	return 8
}
func (s *str) String() string {
	return "string"
}
func (s *str) SExpr() string {
	return "String"
}

var Void = &void{}

type void struct{}

func (v void) Equal(other types.Type) bool {
	return other == Void
}

func (v void) Size() int {
	return 0
}

func (v void) String() string {
	return "void"
}

func (v void) SExpr() string {
	return "Void"
}

type Pointer struct {
	Inner types.Type
}

func (p Pointer) Equal(other types.Type) bool {
	otherPtr, ok := other.(Pointer)
	return ok && p.Inner.Equal(otherPtr.Inner)
}

func (p Pointer) Size() int {
	return 8
}

func (p Pointer) String() string {
	return fmt.Sprintf("%s*", p.Inner)
}

func (p Pointer) SExpr() string {
	return fmt.Sprintf("(Pointer (%s))", p.Inner)
}

var Opaque = &opaque{}

type opaque struct{}

func (o *opaque) Equal(other types.Type) bool {
	return true
}

func (o *opaque) Size() int {
	return 8
}

func (o *opaque) String() string {
	return "opaque"
}

func (o *opaque) SExpr() string {
	return "Opaque"
}
