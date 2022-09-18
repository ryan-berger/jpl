package llvm

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/types"
	"github.com/ryan-berger/jpl/internal/collections"
	"tinygo.org/x/go-llvm"
)

type infixOp func(g *generator, a, b llvm.Value, name string) llvm.Value

func icmpToInfix(predicate llvm.IntPredicate) infixOp {
	return func(g *generator, a, b llvm.Value, name string) llvm.Value {
		return g.builder.CreateICmp(predicate, a, b, name)
	}
}

func fcmpToInfix(predicate llvm.FloatPredicate) infixOp {
	return func(g *generator, a, b llvm.Value, name string) llvm.Value {
		return g.builder.CreateFCmp(predicate, a, b, name)
	}
}

func builderDo(do func(build llvm.Builder, a, b llvm.Value, name string) llvm.Value) infixOp {
	return func(g *generator, a, b llvm.Value, name string) llvm.Value {
		return do(g.builder, a, b, name)
	}
}

func (g *generator) createAssert(cond llvm.Value, msg string) {
	assertFn := g.fns["fail_assertion"]

	failBB := g.ctx.AddBasicBlock(g.curFn.fn, "assertfail")
	contBB := g.ctx.AddBasicBlock(g.curFn.fn, "assertcont")

	g.builder.CreateCondBr(cond, contBB, failBB)

	str := g.builder.CreateGlobalStringPtr(msg, "assert")

	g.builder.SetInsertPointAtEnd(failBB)
	g.builder.CreateCall(assertFn.fn,
		[]llvm.Value{str},
		"",
	)
	g.builder.CreateUnreachable()
	g.builder.SetInsertPointAtEnd(contBB)

}

func (g *generator) div(a, b llvm.Value, name string) llvm.Value {
	div := g.builder.CreateICmp(llvm.IntEQ, b, llvm.ConstInt(g.ctx.Int64Type(), 0, false), "is_zero")

	g.createAssert(g.builder.CreateNot(div, "not_zero"), "division by zero")

	return g.builder.CreateSDiv(a, b, name)
}

func (g *generator) rem(a, b llvm.Value, name string) llvm.Value {
	div := g.builder.CreateICmp(llvm.IntEQ, b, llvm.ConstInt(g.ctx.Int64Type(), 0, false), "is_zero")

	g.createAssert(g.builder.CreateNot(div, "not_zero"), "modulo by zero")

	return g.builder.CreateSRem(a, b, name)
}

var fns = map[types.Type]map[string]infixOp{
	types.Integer: {
		"+":  builderDo(llvm.Builder.CreateAdd),
		"-":  builderDo(llvm.Builder.CreateSub),
		"*":  builderDo(llvm.Builder.CreateMul),
		"/":  (*generator).div,
		"%":  (*generator).rem,
		"<":  icmpToInfix(llvm.IntSLT),
		"<=": icmpToInfix(llvm.IntSLE),
		">":  icmpToInfix(llvm.IntSGT),
		">=": icmpToInfix(llvm.IntSGE),
		"==": icmpToInfix(llvm.IntEQ),
		"!=": icmpToInfix(llvm.IntNE),
	},
	types.Float: {
		"+":  builderDo(llvm.Builder.CreateFAdd),
		"-":  builderDo(llvm.Builder.CreateFSub),
		"*":  builderDo(llvm.Builder.CreateFMul),
		"/":  builderDo(llvm.Builder.CreateFDiv),
		"%":  builderDo(llvm.Builder.CreateFRem),
		"<":  fcmpToInfix(llvm.FloatOLT),
		"<=": fcmpToInfix(llvm.FloatOLE),
		">":  fcmpToInfix(llvm.FloatOGT),
		">=": fcmpToInfix(llvm.FloatOGE),
		"==": fcmpToInfix(llvm.FloatOEQ),
		"!=": fcmpToInfix(llvm.FloatONE),
	},
	types.Boolean: {
		"||": builderDo(llvm.Builder.CreateOr),
		"&&": builderDo(llvm.Builder.CreateAnd),
	},
}

var casts = map[string]func(llvm.Builder, llvm.Context, llvm.Value, string) llvm.Value{
	"int": func(builder llvm.Builder, ctx llvm.Context, value llvm.Value, s string) llvm.Value {
		return builder.CreateFPToSI(value, ctx.Int64Type(), s)
	},
	"float": func(builder llvm.Builder, ctx llvm.Context, value llvm.Value, s string) llvm.Value {
		return builder.CreateSIToFP(value, ctx.DoubleType(), s)
	},
}

func (g *generator) getExpr(val map[string]llvm.Value, expression ast.Expression) llvm.Value {
	switch expr := expression.(type) {
	case *ast.IntExpression:
		return llvm.ConstInt(g.ctx.Int64Type(), uint64(expr.Val), false)
	case *ast.FloatExpression:
		return llvm.ConstFloat(g.ctx.DoubleType(), expr.Val)
	case *ast.BooleanExpression:
		if expr.Val {
			return llvm.ConstInt(g.ctx.Int1Type(), 1, false)
		}
		return llvm.ConstInt(g.ctx.Int1Type(), 0, false)
	case *ast.PrefixExpression:
		r := g.getExpr(val, expr.Expr)
		switch expr.Op {
		case "!":
			return g.builder.CreateNot(r, "not")
		case "-":
			if expr.Expr.Typ() == types.Float {
				return g.builder.CreateFSub(llvm.ConstFloat(g.ctx.DoubleType(), 0.0), r, "neg")
			}
			return g.builder.CreateNeg(r, "neg")
		}
	case *ast.InfixExpression:
		l, r := g.getExpr(val, expr.Left), g.getExpr(val, expr.Right)
		return fns[expr.Left.Typ()][expr.Op](g, l, r, "infx")
	case *ast.CallExpression:
		if fn, ok := casts[expr.Identifier]; ok {
			exp := g.getExpr(val, expr.Arguments[0])
			return fn(g.builder, g.ctx, exp, "cast")
		}

		fun, ok := g.fns[expr.Identifier]
		if !ok {
			panic(fmt.Sprintf("could not find fn %s", expr.Identifier))
		}

		vals := collections.Map(expr.Arguments, func(e ast.Expression) llvm.Value { return g.getExpr(val, e) })

		call := g.builder.CreateCall(fun.fn, make([]llvm.Value, fun.fn.ParamsCount()), "call")
		for i, exp := range vals {
			call.SetOperand(i, exp)
			if ty := exp.Type(); ty.TypeKind() == llvm.PointerTypeKind {
				call.AddCallSiteAttribute(i+1, g.ctx.CreateTypeAttribute(llvm.AttributeKindID("byval"), ty.ElementType()))
			}
		}
		return call
	case *ast.SumTransform:
		return g.genSumTransform(val, expr)
	case *ast.ArrayTransform:
		return g.genArrayTransform(val, expr)
	case *ast.IfExpression:
		return g.genIf(val, expr)
	case *ast.TupleRefExpression:
		return g.builder.CreateExtractValue(g.getExpr(val, expr.Tuple), int(expr.Index), "tuple_val")
	case *ast.ArrayRefExpression:
		arr := g.getExpr(val, expr.Array)
		idxs := collections.Map(expr.Indexes, func(e ast.Expression) llvm.Value { return g.getExpr(val, e) })

		ptr := g.getArrayBase(arr, idxs)

		return g.builder.CreateLoad(ptr, "item")
	case *ast.IdentifierExpression:
		v, ok := val[expr.Identifier]
		if !ok {
			panic(fmt.Sprintf("identifier %s not found", expr.Identifier))
		}
		return v
	case *ast.TupleExpression:
		typ := toLLVMType(g.ctx, expr.Type)
		iv := g.builder.CreateInsertValue(llvm.Undef(typ), g.getExpr(val, expr.Expressions[0]), 0, "first")

		for i := 1; i < len(expr.Expressions); i++ {
			iv = g.builder.CreateInsertValue(iv, g.getExpr(val, expr.Expressions[i]), i, "next")
		}

		return iv
	}
	panic(fmt.Sprintf("unsupported expr: %T", expression))
}
