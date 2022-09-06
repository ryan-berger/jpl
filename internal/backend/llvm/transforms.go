package llvm

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/collections"
	"tinygo.org/x/go-llvm"
)

// initializer will create a phi block with a 0 value for the type t when coming from prev
// If we are coming from the cur block, we need to use the current value v
func (g *generator) initializer(phi, v llvm.Value, prev, cur llvm.BasicBlock) {
	var zero llvm.Value
	switch phi.Type() {
	case g.ctx.Int64Type():
		zero = llvm.ConstInt(phi.Type(), 0, false)
	case g.ctx.FloatType():
		zero = llvm.ConstFloat(phi.Type(), 0)
	}

	phi.AddIncoming(
		[]llvm.Value{
			zero, // 0 for incoming block
			v,    // accumulated value
		},
		[]llvm.BasicBlock{
			prev,
			cur,
		},
	)
}

// recursiveSum generates the following IR pattern:
// ```
// %0:
//
//	%bound = getExpr()
//	%2 = icmp slt i32 0, %bound
//	br i1 %2, label %5, label %3
//
// %2: ; end of the loop
//
//	%4 = phi i64 [ 0, %0 ], [ %11, %7 ]
//	ret i64 %4
//
// %5: ; loop body
//
//	%6 = phi i32 [ %11, %7 ], [ 0, %0 ]
//
// ```
func (g *generator) recursiveSum(vals map[string]llvm.Value, idx int, bindings []ast.OpBinding, exp ast.Expression) (llvm.Value, llvm.BasicBlock) {
	if idx >= len(bindings) {
		return g.getExpr(vals, exp), g.builder.GetInsertBlock()
	}

	b := bindings[idx]

	// get the upper bound of the loop
	bound := g.getExpr(vals, b.Expr)

	// prep the last block
	loopEnd := g.ctx.AddBasicBlock(g.curFn.fn, "loop_end")
	loopBody := g.ctx.AddBasicBlock(g.curFn.fn, "loop_body")

	// make sure to check to see whether we can enter the loop. If we can't, go straight to end
	enterCond := g.builder.CreateICmp(llvm.IntSLT, llvm.ConstInt(g.ctx.Int64Type(), 0, false), bound, "loop_cond_enter")
	g.builder.CreateCondBr(enterCond, loopBody, loopEnd)

	// save the previous basic block for later. If we are entering from prevBB, it means that this is the first time
	// entering. This comes in handy for initializers
	prevBB := g.builder.GetInsertBlock()

	// start generating loop body
	g.builder.SetInsertPointAtEnd(loopBody)

	// create phi node for variable to be incremented
	incPhi := g.builder.CreatePHI(g.ctx.Int64Type(), b.Variable)
	// create phi node for the result of the operations
	resPhi := g.builder.CreatePHI(toLLVMType(g.ctx, exp.Typ()), "res")

	// bind the incremented variable to incPhi and prepare for recursion
	vals[b.Variable] = incPhi

	val, endBB := g.recursiveSum(vals, idx+1, bindings, exp)

	// generate increment expression
	increment := g.builder.CreateAdd(incPhi, llvm.ConstInt(g.ctx.Int64Type(), 1, false), "inc")

	// add up the accumulated total, and the new value
	sum := g.builder.CreateAdd(resPhi, val, "sum")
	g.initializer(resPhi, sum, prevBB, endBB)
	g.initializer(incPhi, increment, prevBB, endBB)

	// check to make sure we can continue looping
	cond := g.builder.CreateICmp(llvm.IntSLT, increment, bound, "loop_cond")
	g.builder.CreateCondBr(cond, loopBody, loopEnd)

	// start generating loop end
	g.builder.SetInsertPointAtEnd(loopEnd)
	totalPhi := g.builder.CreatePHI(resPhi.Type(), "total")
	g.initializer(totalPhi, resPhi, prevBB, endBB)

	return totalPhi, loopEnd
}

func (g *generator) genSumTransform(vals map[string]llvm.Value, t *ast.SumTransform) llvm.Value {
	if len(t.OpBindings) == 0 {
		return g.getExpr(vals, t.Expr)
	}

	cpy := make(map[string]llvm.Value)
	for k, v := range vals {
		cpy[k] = v
	}

	val, _ := g.recursiveSum(cpy, 0, t.OpBindings, t.Expr)
	return val
}

func (g *generator) recursiveArray(vals map[string]llvm.Value, storeTo llvm.Value, idx int, params []llvmBinding) {
	if idx >= len(params) {
		return
	}

	bound := params[idx]

	// prep the last block
	loopEnd := g.ctx.AddBasicBlock(g.curFn.fn, "loop_end")
	loopBody := g.ctx.AddBasicBlock(g.curFn.fn, "loop_body")

	// make sure to check to see whether we can enter the loop. If we can't, go straight to end
	enterCond := g.builder.CreateICmp(llvm.IntSLT, llvm.ConstInt(g.ctx.Int64Type(), 0, false), bound.val, "loop_cond_enter")
	g.builder.CreateCondBr(enterCond, loopBody, loopEnd)

	// save the previous basic block for later. If we are entering from prevBB, it means that this is the first time
	// entering. This comes in handy for initializers
	prevBB := g.builder.GetInsertBlock()

	// start generating loop body
	g.builder.SetInsertPointAtEnd(loopBody)

	// create phi node for variable to be incremented
	incPhi := g.builder.CreatePHI(g.ctx.Int64Type(), "i")

	// bind the incremented variable to incPhi and prepare for recursion
	vals[bound.variable] = incPhi

	nextStore := g.builder.CreateGEP(storeTo, []llvm.Value{incPhi}, "inner_arr")
	g.recursiveArray(vals, nextStore, idx+1, params)

	// generate increment expression
	increment := g.builder.CreateAdd(incPhi, llvm.ConstInt(g.ctx.Int64Type(), 1, false), "inc")
	g.initializer(incPhi, increment, prevBB, loopEnd)

	// check to make sure we can continue looping
	cond := g.builder.CreateICmp(llvm.IntSLT, increment, bound.val, "loop_cond")
	g.builder.CreateCondBr(cond, loopBody, loopEnd)

	// start generating loop end
	g.builder.SetInsertPointAtEnd(loopEnd)
}

type llvmBinding struct {
	variable string
	val      llvm.Value
}

func (g *generator) genArrayTransform(vals map[string]llvm.Value, t *ast.ArrayTransform) llvm.Value {
	if len(t.OpBindings) == 0 {
		return g.getExpr(vals, t.Expr)
	}

	exprs := collections.Map(t.OpBindings, func(o ast.OpBinding) llvm.Value { return g.getExpr(vals, o.Expr) })
	sum := exprs[0]
	for i := 1; i < len(exprs); i++ {
		sum = g.builder.CreateMul(exprs[i], sum, "size_calc")
	}

	var types []llvm.Type
	for i := 0; i < len(exprs); i++ {
		types = append(types, g.ctx.Int64Type())
	}
	expType := toLLVMType(g.ctx, t.Expr.Typ())
	types = append(types, llvm.PointerType(expType, 0))
	structType := g.ctx.StructType(types, false)

	res := g.builder.CreateAlloca(structType, "struct_ptr")

	for i := 0; i < len(exprs); i++ {
		fieldPtr := g.builder.CreateStructGEP(res, i, fmt.Sprintf("rank_%d", i))
		g.builder.CreateStore(exprs[i], fieldPtr)
	}
	arr := g.builder.CreateArrayAlloca(expType, sum, "arr")

	fieldPtr := g.builder.CreateStructGEP(res, len(exprs), "arr_field")
	g.builder.CreateStore(arr, fieldPtr)

	load := g.builder.CreateLoad(res, "load")

	llvmParams := make([]llvmBinding, len(exprs))
	for i := 0; i < len(llvmParams); i++ {
		llvmParams[i] = llvmBinding{
			val:      exprs[i],
			variable: t.OpBindings[i].Variable,
		}
	}

	g.recursiveArray(vals, arr, 0, llvmParams)

	return load
}
