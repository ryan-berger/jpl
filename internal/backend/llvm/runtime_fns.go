package llvm

import (
	"tinygo.org/x/go-llvm"
)

type llvmArg struct {
	name     string
	llvmType llvm.Type
}

type llvmType llvm.Type

type llvmFunc struct {
	name     string
	llvmType llvm.Type
	args     []llvmArg
}

func (g *generator) genRuntime() {
	fnType := llvm.FunctionType(g.ctx.VoidType(),
		[]llvm.Type{llvm.PointerType(g.ctx.Int8Type(), 0)},
		false,
	)

	llvmFun := llvm.AddFunction(g.module, "fail_assertion", fnType)
	llvmFun.SetLinkage(llvm.ExternalLinkage)
	llvmFun.Param(0).SetName("message")

	g.fns["fail_assertion"] = fn{
		fn: llvmFun,
		params: map[string]llvm.Value{
			"message": llvmFun.Param(0),
		},
	}
}
