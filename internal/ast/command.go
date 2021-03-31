package ast

import (
	"bytes"
	"fmt"

	"github.com/ryan-berger/jpl/internal/types"
)

type Command interface {
	Node
	command()
}

type Function struct {
	Var        string
	Bindings   []Binding
	ReturnType types.Type
	Statements []Statement
	Location
}

func (f *Function) SExpr() string {
	bindings := make([]string, len(f.Bindings))
	stmts := make([]string, len(f.Statements))

	for i, stmt := range f.Statements {
		stmts[i] = stmt.SExpr()
	}

	for i, b := range f.Bindings {
		bindings[i] = b.SExpr()
	}

	return fmt.Sprintf("(Func %s %s)", f.Var, f.ReturnType)
}

func (f *Function) command() {}
func (f *Function) String() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("fn %s (", f.Var))
	for i, b := range f.Bindings {
		buf.WriteString(b.String())
		if i != len(f.Bindings)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(fmt.Sprintf(") : %s {\n", f.ReturnType))
	for _, s := range f.Statements {
		buf.WriteString(fmt.Sprintf("\t%s\n", s))
	}

	buf.WriteString("}\n")
	return buf.String()
}

type Read struct {
	Type     string
	Src      string
	Argument Argument
	Location
}

func (r *Read) SExpr() string {
	return fmt.Sprintf("(ReadImageCmd %s %s)", r.Src, r.Argument.SExpr())
}

func (r *Read) String() string {
	return fmt.Sprintf("read %s %s to %s", r.Type, r.Src, r.Argument)
}
func (r *Read) command() {}

type Write struct {
	Type string
	Expr Expression
	Dest string
	Location
}

func (w *Write) SExpr() string {
	return fmt.Sprintf("(WriteImageCmd %s %s)", w.Expr.SExpr(), w.Dest)
}

func (w *Write) String() string {
	return fmt.Sprintf("write %s %s to %s", w.Type, w.Expr, w.Dest)
}
func (w *Write) command() {}

type Show struct {
	Expr Expression
	Location
}

func (s *Show) SExpr() string {
	return fmt.Sprintf("(ShowCmd %s)", s.Expr.SExpr())
}

func (s *Show) String() string {
	return fmt.Sprintf("show %s", s.Expr)
}
func (s *Show) command() {}

type Print struct {
	Str string
	Location
}

func (p *Print) SExpr() string {
	return fmt.Sprintf("(PrintCmd %s)", p.Str)
}

func (p *Print) String() string {
	return fmt.Sprintf("print %s", p.Str)
}

func (p *Print) command() {}

type Time struct {
	Command Command
	Location
}

func (t *Time) SExpr() string {
	cmd := t.Command.SExpr()
	if _, ok := t.Command.(Statement); ok {
		cmd = fmt.Sprintf("(StmtCmd %s)", cmd)
	}
	return fmt.Sprintf("(TimeCmd %s)", cmd)
}

func (t *Time) String() string {
	return fmt.Sprintf("time %s", t.Command)
}
func (t *Time) command() {}
