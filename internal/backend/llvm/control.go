package llvm

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"tinygo.org/x/go-llvm"
)

func (g *generator) genIf(vals map[string]llvm.Value, expr *ast.IfExpression) llvm.Value {
	thenBB := g.ctx.AddBasicBlock(g.curFn.fn, "then")
	elseBB := g.ctx.AddBasicBlock(g.curFn.fn, "else")
	contBB := g.ctx.AddBasicBlock(g.curFn.fn, "ifcont")

	// if cond: goto then else goto
	cond := g.getExpr(vals, expr.Condition)
	g.builder.CreateCondBr(cond, thenBB, elseBB)

	// then:
	g.builder.SetInsertPointAtEnd(thenBB)
	cons := g.getExpr(vals, expr.Consequence)
	g.builder.CreateBr(contBB)

	thenBB = g.builder.GetInsertBlock()

	// else:
	g.builder.SetInsertPointAtEnd(elseBB)
	other := g.getExpr(vals, expr.Otherwise)
	g.builder.CreateBr(contBB)

	elseBB = g.builder.GetInsertBlock()

	g.builder.SetInsertPointAtEnd(contBB)
	// combine
	phi := g.builder.CreatePHI(toLLVMType(g.ctx, expr.Type), "phi")
	phi.AddIncoming([]llvm.Value{cons, other}, []llvm.BasicBlock{thenBB, elseBB})
	return phi
}
