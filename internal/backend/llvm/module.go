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
	curFn   fn
	fns     map[string]fn
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
		fns:     make(map[string]fn),
	}

	fpm := llvm.NewFunctionPassManagerForModule(module)
	fpm.AddCFGSimplificationPass()
	fpm.AddReassociatePass()
	fpm.AddInstructionCombiningPass()
	fpm.InitializeFunc()

	defer fpm.Dispose()

	g.genRuntime()
	g.generate(p, fpm)

	module.Dump()

}

func (g *generator) generate(p ast.Program, fpm llvm.PassManager) {
	for _, cmd := range p {
		if fn, ok := cmd.(*ast.Function); ok {
			g.declareFunction(fn)
		}
	}

	for _, cmd := range p {
		if fn, ok := cmd.(*ast.Function); ok {
			g.genFunction(fn)
			g.fns[fn.Var].fn.Dump()
			//fpm.RunFunc(g.fns[fn.Var].fn)
		}
	}
}

func curryType(ctx llvm.Context) func(p types.Type) llvm.Type {
	return func(p types.Type) llvm.Type {
		return toLLVMType(ctx, p)
	}
}

func toLLVMType(ctx llvm.Context, p types.Type) llvm.Type {
	switch {
	case p == types.Float:
		return ctx.FloatType()
	case p == types.Integer:
		return ctx.Int64Type()
	case p == types.Boolean:
		return ctx.Int1Type()
	}

	switch t := p.(type) {
	case *types.Array:
		ranks := make([]llvm.Type, t.Rank)
		for i := 0; i < t.Rank; i++ {
			ranks[i] = ctx.Int64Type()
		}

		inner := toLLVMType(ctx, t.Inner)
		inner = llvm.PointerType(inner, 0)

		return ctx.StructType(append(ranks, inner), false)
	case *types.Tuple:
		return ctx.StructType(collections.Map(t.Types, curryType(ctx)), false)
	}

	panic("unreachable")
}

type infixOp func(builder llvm.Builder, a, b llvm.Value, name string) llvm.Value

func icmpToInfix(predicate llvm.IntPredicate) infixOp {
	return func(builder llvm.Builder, a, b llvm.Value, name string) llvm.Value {
		return builder.CreateICmp(predicate, a, b, name)
	}
}

func fcmpToInfix(predicate llvm.FloatPredicate) infixOp {
	return func(builder llvm.Builder, a, b llvm.Value, name string) llvm.Value {
		return builder.CreateFCmp(predicate, a, b, name)
	}
}

var fns = map[types.Type]map[string]infixOp{
	types.Integer: {
		"+":  llvm.Builder.CreateAdd,
		"-":  llvm.Builder.CreateSub,
		"*":  llvm.Builder.CreateMul,
		"/":  llvm.Builder.CreateSDiv,
		"<":  icmpToInfix(llvm.IntSLT),
		"<=": icmpToInfix(llvm.IntSLE),
		">":  icmpToInfix(llvm.IntSGT),
		">=": icmpToInfix(llvm.IntSGE),
		"==": icmpToInfix(llvm.IntEQ),
		"!=": icmpToInfix(llvm.IntNE),
	},
	types.Float: {
		"+":  llvm.Builder.CreateFAdd,
		"-":  llvm.Builder.CreateFSub,
		"*":  llvm.Builder.CreateFMul,
		"/":  llvm.Builder.CreateFDiv,
		"<":  fcmpToInfix(llvm.FloatOLT),
		"<=": fcmpToInfix(llvm.FloatOLE),
		">":  fcmpToInfix(llvm.FloatOGT),
		">=": fcmpToInfix(llvm.FloatOGE),
		"==": fcmpToInfix(llvm.FloatOEQ),
		"!=": fcmpToInfix(llvm.FloatONE),
	},
	types.Boolean: {
		"||": llvm.Builder.CreateOr,
		"&&": llvm.Builder.CreateAnd,
	},
}

func (g *generator) getExpr(val map[string]llvm.Value, expression ast.Expression) llvm.Value {
	switch expr := expression.(type) {
	case *ast.IntExpression:
		return llvm.ConstInt(g.ctx.Int64Type(), uint64(expr.Val), false)
	case *ast.FloatExpression:
		return llvm.ConstFloat(g.ctx.FloatType(), expr.Val)
	case *ast.BooleanExpression:
		if expr.Val {
			return llvm.ConstInt(g.ctx.Int1Type(), 1, false)
		}
		return llvm.ConstInt(g.ctx.Int1Type(), 0, false)
	case *ast.PrefixExpression:
		r := g.getExpr(val, expr.Expr)
		switch expr.Op {
		case "-":
			return g.builder.CreateNeg(r, "neg")
		}
	case *ast.InfixExpression:
		l, r := g.getExpr(val, expr.Left), g.getExpr(val, expr.Right)
		return fns[expr.Left.Typ()][expr.Op](g.builder, l, r, "infx")
	case *ast.CallExpression:
		fun, ok := g.fns[expr.Identifier]
		if !ok {
			panic(fmt.Sprintf("call expr"))
		}
		args := make([]llvm.Value, len(fun.params))
		for i, e := range expr.Arguments {
			args[i] = g.getExpr(val, e)
		}
		return g.builder.CreateCall(fun.fn, args, "call")
	case *ast.SumTransform:
		return g.genSumTransform(val, expr)
	case *ast.ArrayTransform:
		return g.genArrayTransform(val, expr)
	case *ast.IfExpression:
		return g.genIf(val, expr)
	case *ast.TupleRefExpression:
		return g.builder.CreateExtractValue(g.getExpr(val, expr.Tuple), int(expr.Index), "tuple_lookup")
	case *ast.ArrayRefExpression:
		arr := g.getExpr(val, expr.Array)

		for _, idx := range expr.Indexes {
			arr = g.builder.CreateGEP(arr, []llvm.Value{g.getExpr(val, idx)}, "elem_ptr")
		}

		return g.builder.CreateLoad(arr, "elem")
	case *ast.IdentifierExpression:
		v, ok := val[expr.Identifier]
		if !ok {
			panic(fmt.Sprintf("identifier %s not found", expr.Identifier))
		}
		return v
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
	case *ast.AssertStatement:
		assertFn := g.fns["fail_assertion"]
		exp := g.getExpr(m, stmt.Expr)

		failBB := g.ctx.AddBasicBlock(g.curFn.fn, "assertfail")
		contBB := g.ctx.AddBasicBlock(g.curFn.fn, "assertcont")

		g.builder.CreateCondBr(exp, failBB, contBB)

		str := g.builder.CreateGlobalStringPtr(stmt.Message, "assert")

		g.builder.SetInsertPointAtEnd(failBB)
		g.builder.CreateCall(assertFn.fn,
			[]llvm.Value{str},
			"",
		)
		g.builder.CreateUnreachable()
		g.builder.SetInsertPointAtEnd(contBB)
	default:
		panic(fmt.Sprintf("unsupported stmt: %T", s))
	}
}
