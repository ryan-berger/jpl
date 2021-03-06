package ast

import (
	"bytes"
	"fmt"

	"github.com/ryan-berger/jpl/internal/meta"
)

type SExpr interface {
	SExpr() string
}

type Node interface {
	SExpr
	meta.Locationer
	String() string
}

type Program []Command

func (p Program) SExpr() string {
	buf := bytes.NewBufferString("")
	for _, cmd := range p {
		if _, ok := cmd.(Statement); ok {
			buf.WriteString(fmt.Sprintf("(StmtCmd %s)", cmd.SExpr()))
			continue
		}
		buf.WriteString(cmd.SExpr())
	}
	return buf.String()
}

func (p Program) Loc() (int, int) {
	return 0, 0
}

func (p Program) String() string {
	return ""
}
