package generator

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/types"
)

func ident(let *ast.LetStatement) string {
	arg, ok := let.LValue.(*ast.VariableArgument)
	if !ok {
		panic("must call ident() on VariableArgument")
	}
	return arg.Variable
}

const assertAsm = `cmp dword [rbp - %d], 0
jne %[2]s
lea rdi, [rel %s]
call _fail_assertion
%[2]s

`

const printAsm = `lea rdi, [rel %s]
call _print

`

func (g *generator) genCommand(command ast.Command) {

	switch cmd := command.(type) {
	case *ast.AssertStatement:
		loc := g.frame[cmd.Expr.String()]
		msg := g.mapper[cmd]
		g.buf.WriteString(fmt.Sprintf(assertAsm, loc, ".SKIP", msg))
	case *ast.Print:
		msg := g.mapper[cmd]
		g.buf.WriteString(fmt.Sprintf(printAsm, msg))
	case ast.Statement:
		g.genStatement(cmd)
	}
}

const returnInt = `mov rax, [rbp - %d]
add rsp, %d
pop rbp
ret
`

func (g *generator) genStatement(
	statement ast.Statement) {
	switch stmt := statement.(type) {
	case *ast.LetStatement:
		g.genLetStatement(stmt, g.frame)
	case *ast.ReturnStatement:
		ident := stmt.Expr.(*ast.IdentifierExpression).Identifier
		g.buf.WriteString(fmt.Sprintf(returnInt, g.frame[ident], g.size))
	}
}

const letConstant = `mov rbx, [rel %s]
mov [rbp - %d], rbx

`

const moveNumber = `mov rbx [rbp - %d]
mov [rbp - %d], rbp

`

const moveBool = `mov ebx, [rbp - %d]
mov [rbp - %d], ebx

`

func (g *generator) genLetStatement(
	let *ast.LetStatement,
	f frame) {

	loc := f[ident(let)]
	switch exp := let.Expr.(type) {
	case *ast.FloatExpression, *ast.IntExpression:
		constName := g.mapper[let]

		g.buf.WriteString(fmt.Sprintf(letConstant, constName, loc))
	case *ast.IdentifierExpression:
		locR := f[exp.Identifier]

		switch exp.Type {
		case types.Integer, types.Float:
			g.buf.WriteString(fmt.Sprintf(moveNumber, loc, locR))
		case types.Boolean:
			g.buf.WriteString(fmt.Sprintf(moveBool, loc, locR))
		}
	case *ast.CallExpression:
		g.callExpressionPlanner(loc, exp, f)
	}
}

var intRegisters = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}
var floatRegisters = []string{
	"xmm0", "xmm1", "xmm2", "xmm3", "xmm4", "xmm5", "xmm6", "xmm7",
}

var intArg = `mov %s, [rbp - %d]
`

var floatArg = `movsd %s, [rbp - %d]
`

func (g *generator) callExpressionPlanner(
	retLoc int, expr *ast.CallExpression, f frame) {
	
	intReg := 0   // keep track of number of int/bool args used
	floatReg := 0 // keep track of number of float args used

	// stackSize := 0
	// stackArg := bytes.NewBuffer([]byte{})

	for _, arg := range expr.Arguments {
		ref := arg.(*ast.IdentifierExpression)
		loc := f[ref.Identifier]
		switch typ := arg.Typ(); typ {
		case types.Integer:
			g.buf.WriteString(fmt.Sprintf(intArg, intRegisters[intReg], loc))
			intReg++
		case types.Boolean:
			intReg++
		case types.Float:
			g.buf.WriteString(fmt.Sprintf(floatArg, floatRegisters[floatReg], loc))
			floatReg++
		default:
			arr, ok := typ.(*types.Array)
			if ok {
				arr.Size()
				continue
			}

		}
	}

	g.buf.WriteString(fmt.Sprintf("call _%s\n", expr.Identifier))
	switch expr.Type {
	case types.Integer:
		g.buf.WriteString(fmt.Sprintf("mov [rbp - %d], rax\n\n", retLoc))
	case types.Boolean:
		g.buf.WriteString(fmt.Sprintf("movsd [rbp - %d], eax\n\n", retLoc))
	case types.Float:
		g.buf.WriteString(fmt.Sprintf("movsd [rbp - %d], xmm0\n\n", retLoc))
	}
}
