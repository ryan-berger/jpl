package ast

import (
	"fmt"
	"strings"
)

type Argument interface {
	SExpr
	String() string
	argument()
	lValue()
}

type VariableArgument struct {
	Variable string
}

func (v *VariableArgument) SExpr() string {
	return fmt.Sprintf("(VarArgument %s)", v.Variable)
}

func (v *VariableArgument) String() string { return v.Variable }
func (v *VariableArgument) argument()      {}
func (v *VariableArgument) lValue()        {}

type VariableArr struct {
	Variable  string
	Variables []string
}
// TODO: make this work when needed
func (v *VariableArr) SExpr() string { return "" }

func (v *VariableArr) String() string {
	return fmt.Sprintf("%s[%s]",
		v.Variable, strings.Join(v.Variables, ", "))
}
func (v *VariableArr) argument() {}
func (v *VariableArr) lValue()   {}
