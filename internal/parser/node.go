package parser

import (
	"bytes"
	"fmt"
	"strings"
)

type Program struct {
	Commands []Command
}

// Bindings:

type VarBinding struct {
	Variable string
	Type     string
}

func (b *VarBinding) String() string {
	return fmt.Sprintf("%s : %s", b.Variable, b.Type)
}
func (b *VarBinding) binding() {}

type TupleBinding []Binding
func (b TupleBinding) String() string {
	strs := make([]string, len(b))
	for i, b := range b {
		strs[i] = b.String()
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ","))
}
func (b TupleBinding) binding() {}

// cmd:

type Function struct {
	Var         string
	Bindings    []Binding
	ReturnType  string
	Expressions []Expression
}

func (f *Function) TokenLiteral() string { return "fn" }

func (f *Function) String() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("%s (", f.Var))
	for i, b := range f.Bindings {
		buf.WriteString(b.String())
		if i != len(f.Bindings)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString(")\n {\n")
	for _, s := range f.Expressions {
		buf.WriteString(s.String())
	}

	buf.WriteString("}\n")
	return buf.String()
}

type LetStatement struct {
	LValue Binding
	Expr   Expression
}
