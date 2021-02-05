package ast

import (
	"bytes"
	"fmt"
)

type Command interface {
	String() string
	command()
}

type Function struct {
	Var        string
	Bindings   []Binding
	ReturnType string
	Statements []Statement
}

func (f *Function) command() {}
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
	for _, s := range f.Statements {
		buf.WriteString(s.String())
	}

	buf.WriteString("}\n")
	return buf.String()
}
