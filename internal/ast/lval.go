package ast

import (
	"fmt"
	"strings"
)

// LValue is the interface for an lvalue of the form:
// lvalue : <argument>
//        | { <lvalue> , ... }
type LValue interface {
	Node
	// Type() Type
	lValue()
}

// LTuple is the struct that contains the lvalue production {<lvalue> , ...}
type LTuple struct {
	Args []LValue
	Location
}

// SExpr is an implementation of the SExpr interface
func (l *LTuple) SExpr() string {
	panic("implement me")
}

func (l *LTuple) String() string {
	lVals := make([]string, len(l.Args))
	for i := 0; i < len(l.Args); i++ {
		lVals[i] = l.Args[i].String()
	}
	return fmt.Sprintf("{%s}", strings.Join(lVals, ", "))
}
func (l *LTuple) lValue() {}
