package ast

import (
	"fmt"
	"strings"

	"github.com/ryan-berger/jpl/internal/types"
)

type Argument interface {
	LValue
	argument()
}

type VariableArgument struct {
	Variable string
	Type     types.Type
	Location
}

func (v *VariableArgument) SExpr() string {
	if v.Type != nil {
		return fmt.Sprintf("(VarArgument %s %s)", v.Type, v.Variable)
	}
	return fmt.Sprintf("(VarArgument %s)", v.Variable)
}

func (v *VariableArgument) String() string { return v.Variable }
func (v *VariableArgument) argument()      {}
func (v *VariableArgument) lValue()        {}

type VariableArr struct {
	Variable  string
	Variables []string
	Type      types.Type
	Location
}

// TODO: make this work when needed
func (v *VariableArr) SExpr() string {
	sExps := make([]string, len(v.Variables))
	for i, v := range v.Variables {
		sExps[i] = fmt.Sprintf("(VarArgument %s)", v)
	}
	return fmt.Sprintf("(ArrayVar %s %s)", v.Variable, strings.Join(sExps, " "))
}

func (v *VariableArr) String() string {
	return fmt.Sprintf("%s[%s]",
		v.Variable, strings.Join(v.Variables, ", "))
}
func (v *VariableArr) argument() {}
func (v *VariableArr) lValue()   {}
