package ast

import (
	"bytes"
	"fmt"
)

type Command interface {
	SExpr
	String() string
	command()
}

type Function struct {
	Var        string
	Bindings   []Binding
	ReturnType Type
	Statements []Statement
}

func (f *Function) SExpr() string {
	panic("implement me")
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
}

func (r *Read) SExpr() string {
	return fmt.Sprintf("(ReadImageCmd %s %s)", r.Argument.SExpr(), r.Src)
}

func (r *Read) String() string {
	return fmt.Sprintf("read %s %s to %s", r.Type, r.Src, r.Argument)
}
func (r *Read) command() {}

type Write struct {
	Type string
	Expr Expression
	Dest string
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
}

func (t *Time) SExpr() string {
	return fmt.Sprintf("(TimeCmd %s)", t.Command.SExpr())
}

func (t *Time) String() string {
	return fmt.Sprintf("time %s", t.Command)
}
func (t *Time) command() {}
