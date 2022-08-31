package llvm

import (
	"fmt"
	"io"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/types"
	"github.com/ryan-berger/jpl/internal/collections"
	"github.com/ryan-berger/jpl/internal/symbol"
	"tinygo.org/x/go-llvm"
)

type generator struct {
	ctx     llvm.Context
	builder llvm.Builder
	module  llvm.Module
}

func Generate(p ast.Program, s *symbol.Table, w io.Writer) {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	ctx := llvm.NewContext()
	module := ctx.NewModule("main")
	builder := ctx.NewBuilder()

	defer builder.Dispose()

	g := generator{
		ctx:     ctx,
		builder: builder,
		module:  module,
	}
	fmt.Println("hi")
	g.generate(p)

	fmt.Printf("module: %s\n", module.String())

}

func (g *generator) generate(p ast.Program) {
	for _, cmd := range p {
		if fn, ok := cmd.(*ast.Function); ok {
			g.genFunction(fn)
		}
	}
}

func toLLVMType(p types.Type) llvm.Type {
	switch {
	case p == types.Float:
		return llvm.FloatType()
	case p == types.Integer:
		return llvm.Int64Type()
	case p == types.Boolean:
		return llvm.Int1Type()
	}

	switch t := p.(type) {
	case *types.Array:
		return llvm.ArrayType(toLLVMType(p), t.Rank)
	case *types.Tuple:
		return llvm.StructType(collections.Map(t.Types, toLLVMType), false)
	}

	panic("unreachable")
}

func bindingToLLVMType(b ast.Binding) llvm.Type {
	switch bind := b.(type) {
	case *ast.TypeBind:
		return toLLVMType(bind.Type)
	case *ast.TupleBinding:
		return llvm.StructType(collections.Map(bind.Bindings, bindingToLLVMType), false)
	default:
		panic("unreachable")
	}
}

func (g *generator) genFunction(f *ast.Function) {
	fnType := llvm.FunctionType(
		toLLVMType(f.ReturnType),
		collections.Map(f.Bindings, bindingToLLVMType), false)

	fn := llvm.AddFunction(g.module, f.Var, fnType)

	m := make(map[string]llvm.Value)

	for i, f := range f.Bindings {
		bind := f.(*ast.TypeBind)
		variable := bind.Argument.(*ast.Variable)

		fn.Param(i).SetName(variable.Variable)
		m[variable.Variable] = fn.Param(i)
	}

	bb := g.ctx.AddBasicBlock(fn, fmt.Sprintf("%s_bb", f.Var))
	g.builder.SetInsertPointAtEnd(bb)

	for _, s := range f.Statements {
		g.generateStatement(m, s)
	}
}

func (g *generator) getExpr(val map[string]llvm.Value, expression ast.Expression) llvm.Value {
	switch expr := expression.(type) {
	case *ast.IntExpression:
		return llvm.ConstInt(llvm.Int64Type(), uint64(expr.Val), false)
	case *ast.FloatExpression:
		return llvm.ConstFloat(llvm.FloatType(), expr.Val)
	case *ast.PrefixExpression:
		r := g.getExpr(val, expr.Expr)
		switch expr.Op {
		case "-":
			return g.builder.CreateNeg(r, "neg")
		}
	case *ast.InfixExpression:
		l, r := g.getExpr(val, expr.Left), g.getExpr(val, expr.Right)
		switch expr.Op {
		case "+":
			return g.builder.CreateAdd(l, r, "add")
		case "-":
			return g.builder.CreateSub(l, r, "sub")
		case "*":
			return g.builder.CreateMul(l, r, "mul")
		case "/":
			return g.builder.CreateSDiv(l, r, "div")
		}
	case *ast.IdentifierExpression:
		return val[expr.Identifier]
	}

	panic(fmt.Sprintf("unsupported expr: %T", expression))
}

func (g *generator) generateStatement(m map[string]llvm.Value, s ast.Statement) {
	switch stmt := s.(type) {
	case *ast.LetStatement:
		l := stmt.LValue.(*ast.Variable)
		exp := g.getExpr(m, stmt.Expr)
		m[l.Variable] = exp
	case *ast.ReturnStatement:
		g.builder.CreateRet(g.getExpr(m, stmt.Expr))
	default:
		panic(fmt.Sprintf("unsupported stmt: %T", s))
	}
}
