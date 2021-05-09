package ast

import (
	"fmt"
)

// Statement is the abstract interface for a statement of the form:
// stmt : let <lvalue> = <expr>
//      | assert <expr> , <string>
//      | return <expr>
type Statement interface {
	Command
	statement()
}


// LetStatement holds an LValue declaration, and the assignment expression Expr for
// the statement production:
// let <lvalue> = <expr>
type LetStatement struct {
	LValue LValue
	Expr   Expression
	Location
}

// SExpr is an implementation of the SExpr interface
func (l *LetStatement) SExpr() string {
	return fmt.Sprintf("(LetStmt (ArgLValue %s) %s)", l.LValue.SExpr(), l.Expr.SExpr())
}

func (l *LetStatement) command()   {}
func (l *LetStatement) statement() {}

func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", l.LValue, l.Expr)
}

// ReturnStatement holds an expression Expr and represents the production:
// return <expression>
type ReturnStatement struct {
	Expr Expression
	Location
}

// SExpr is an implementation of the SExpr interface
func (r *ReturnStatement) SExpr() string {
	return fmt.Sprintf("(ReturnStmt %s)", r.Expr.SExpr())
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", r.Expr.String())
}
func (r *ReturnStatement) command()   {}
func (r *ReturnStatement) statement() {}


// AssertStatement holds an expression Expr and message to throw if Expr does not evaluate to true.
// It represents the production:
// assert <expression> , <string>
type AssertStatement struct {
	Expr    Expression
	Message string
	Location
}

// SExpr is an implementation of the SExpr interface
func (a *AssertStatement) SExpr() string {
	return fmt.Sprintf("(AssertStmt %s %s)", a.Expr.SExpr(), a.Message)
}

func (a *AssertStatement) String() string {
	return fmt.Sprintf("assert %s , %s", a.Expr.String(), a.Message)
}

func (a *AssertStatement) command()   {}
func (a *AssertStatement) statement() {}


type AttributeStatement struct {
	Annotation string
}

func (a *AttributeStatement) SExpr() string {
	return fmt.Sprintf("(Attribute %s)", a.Annotation)
}

func (a AttributeStatement) Loc() (int, int) {
	return 0, 0
}

func (a AttributeStatement) String() string {
	return fmt.Sprintf("attribute %s", a.Annotation)
}

func (a AttributeStatement) command() {}
func (a AttributeStatement) statement() {}
