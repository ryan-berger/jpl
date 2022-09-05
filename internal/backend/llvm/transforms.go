package llvm

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"tinygo.org/x/go-llvm"
)

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

	enterCond := g.builder.CreateICmp(llvm.IntSLT, llvm.ConstInt(g.ctx.Int64Type(), 0, false), bound, "loop_cond_enter")
	g.builder.CreateCondBr(enterCond, loopBody, loopEnd)

	prevBB := g.builder.GetInsertBlock()

	g.builder.SetInsertPointAtEnd(loopBody)

	incPhi := g.builder.CreatePHI(g.ctx.Int64Type(), b.Variable)
	resPhi := g.builder.CreatePHI(toLLVMType(g.ctx, exp.Typ()), "res")

	vals[b.Variable] = incPhi

	val, endBB := g.recursiveSum(vals, idx+1, bindings, exp)

	increment := g.builder.CreateAdd(incPhi, llvm.ConstInt(g.ctx.Int64Type(), 1, false), "inc")
	sum := g.builder.CreateAdd(resPhi, val, "sum")
	resPhi.AddIncoming(
		[]llvm.Value{
			llvm.ConstInt(resPhi.Type(), 0, false),
			sum,
		},
		[]llvm.BasicBlock{
			prevBB,
			endBB,
		},
	)

	incPhi.AddIncoming(
		[]llvm.Value{
			llvm.ConstInt(g.ctx.Int64Type(), 0, false),
			increment,
		},
		[]llvm.BasicBlock{
			prevBB,
			endBB,
		},
	)

	cond := g.builder.CreateICmp(llvm.IntSLT, increment, bound, "loop_cond")
	g.builder.CreateCondBr(cond, loopBody, loopEnd)

	g.builder.SetInsertPointAtEnd(loopEnd)

	total := g.builder.CreatePHI(resPhi.Type(), "total")
	total.AddIncoming([]llvm.Value{
		llvm.ConstInt(resPhi.Type(), 0, true),
		resPhi,
	}, []llvm.BasicBlock{
		prevBB,
		endBB,
	})

	return total, loopEnd
}

func (g *generator) genSumTransform(vals map[string]llvm.Value, t *ast.SumTransform) llvm.Value {
	val, _ := g.recursiveSum(vals, 0, t.OpBindings, t.Expr)
	return val
}
