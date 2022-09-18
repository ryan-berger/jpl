package llvm

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast/types"
	runtime2 "github.com/ryan-berger/jpl/internal/backend/llvm/runtime"
	"github.com/ryan-berger/jpl/internal/collections"
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

func getRuntimeType(ctx llvm.Context, t types.Type) llvm.Type {
	switch t {
	case runtime2.Void:
		return ctx.VoidType()
	case runtime2.String:
		return llvm.PointerType(ctx.Int8Type(), 0)
	case runtime2.Opaque:
		return ctx.Int8Type()
	default:
		ptr, ok := t.(*runtime2.Pointer)
		if ok {
			fmt.Printf("pointer type: %s\n", ptr.Inner.String())
			return llvm.PointerType(getRuntimeType(ctx, ptr.Inner), 0)
		}

		return toLLVMType(ctx, t)
	}
}

func (g *generator) genRuntimeFn(runtimeFn runtime2.Function) fn {
	ret := getRuntimeType(g.ctx, runtimeFn.Return)
	args := collections.Map(runtimeFn.Params, func(p runtime2.Param) llvm.Type {
		return getRuntimeType(g.ctx, p.Type)
	})

	fnType := llvm.FunctionType(ret, args, false)

	llvmFn := llvm.AddFunction(g.module, "_"+runtimeFn.Name, fnType)
	llvmFn.SetLinkage(llvm.ExternalLinkage)

	f := fn{
		fn:     llvmFn,
		params: make(map[string]llvm.Value),
	}

	for i, p := range runtimeFn.Params {
		llvmFn.Param(i).SetName(p.Name)
		f.params[p.Name] = llvmFn.Param(i)
	}

	llvmFn.Dump()

	return f
}

func (g *generator) genRuntime() {
	for k, v := range runtime2.Functions {
		g.fns[k] = g.genRuntimeFn(v)
	}
}
