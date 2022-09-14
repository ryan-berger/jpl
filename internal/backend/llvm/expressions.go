package llvm

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/types"
	"github.com/ryan-berger/jpl/internal/collections"
	"tinygo.org/x/go-llvm"
)

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
			return g.builder.CreateNeg(r, "neg")
		}
	case *ast.InfixExpression:
		l, r := g.getExpr(val, expr.Left), g.getExpr(val, expr.Right)
		return fns[expr.Left.Typ()][expr.Op](g.builder, l, r, "infx")
	case *ast.CallExpression:
		fun, ok := g.fns[expr.Identifier]
		if !ok {
			panic(fmt.Sprintf("could not find fn %s", expr.Identifier))
		}

		call := g.builder.CreateCall(fun.fn, make([]llvm.Value, fun.fn.ParamsCount()), "call")
		for i, e := range expr.Arguments {
			exp := g.getExpr(val, e)
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
