package ast

import (
	"fmt"
	"strings"
)

type Argument interface {
	String() string
	argument()
	lValue()
}

type VariableArgument struct {
	Variable string
}

func (v *VariableArgument) String() string { return v.Variable }
func (v *VariableArgument) argument()      {}
func (v *VariableArgument) lValue()        {}

type VariableArr struct {
	Variable  string
	Variables []string
}

func (v *VariableArr) String() string {
	return fmt.Sprintf("%s[%s]",
		v.Variable, strings.Join(v.Variables, ", "))
}
func (v *VariableArr) argument() {}
func (v *VariableArr) lValue()   {}
