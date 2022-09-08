package llvm

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/collections"
	"tinygo.org/x/go-llvm"
)

type fn struct {
	fn     llvm.Value
	params map[string]llvm.Value
}

func curryBinding(ctx llvm.Context) func(b ast.Binding) llvm.Type {
	return func(b ast.Binding) llvm.Type {
		return bindingToLLVMType(ctx, b)
	}
}
func bindingToLLVMType(ctx llvm.Context, b ast.Binding) llvm.Type {
	switch bind := b.(type) {
	case *ast.TypeBind:
		return toLLVMType(ctx, bind.Type)
	case *ast.TupleBinding:
		return ctx.StructType(collections.Map(bind.Bindings, curryBinding(ctx)), false)
	default:
		panic("unreachable")
	}
}

func (g *generator) declareFunction(f *ast.Function) {
	fnType := llvm.FunctionType(
		toLLVMType(g.ctx, f.ReturnType),
		collections.Map(f.Bindings, curryBinding(g.ctx)), false)

	llvmFn := llvm.AddFunction(g.module, f.Var, fnType)

	m := make(map[string]llvm.Value)

	for i, f := range f.Bindings {
		switch bind := f.(type) {
		case *ast.TypeBind:
			switch arg := bind.Argument.(type) {
			case *ast.Variable:
				llvmFn.Param(i).SetName(arg.Variable)
				m[arg.Variable] = llvmFn.Param(i)
			case *ast.VariableArr:
				llvmFn.Param(i).SetName(arg.Variable)
			}
		case *ast.TupleBinding:
			llvmFn.Param(i).SetName(fmt.Sprintf("struct_%d", i))
		}
	}

	g.fns[f.Var] = fn{
		fn:     llvmFn,
		params: m,
	}
}

func (g *generator) makeTupBinding(base llvm.Value, b ast.Binding, idxs ...int) {
	switch bind := b.(type) {
	case *ast.TypeBind:
		if arg, ok := bind.Argument.(*ast.Variable); ok {
			val := base
			for _, arg := range idxs {
				val = g.builder.CreateExtractValue(val, arg, "get")
			}
			g.curFn.params[arg.Variable] = val
		}
	case *ast.TupleBinding:
		for i, innerBind := range bind.Bindings {
			g.makeTupBinding(base, innerBind, append(idxs, i)...)
		}
	}
}

func (g *generator) genBindingAccesses(base llvm.Value, b ast.Binding) {
	switch bind := b.(type) {
	case *ast.TypeBind:
		switch arg := bind.Argument.(type) {
		case *ast.Variable:
			return
		case *ast.VariableArr:
			g.curFn.params[arg.Variable] = base

			for i, v := range arg.Variables {
				g.curFn.params[v] = g.builder.CreateExtractValue(base, i, "get_rank")
			}
		}
	case *ast.TupleBinding:
		for i, tupBind := range bind.Bindings {
			g.makeTupBinding(base, tupBind, i)
		}
	}
}

func (g *generator) genFunction(f *ast.Function) {
	fun, ok := g.fns[f.Var]
	if !ok {
		panic(fmt.Sprintf("fn %s not found", f.Var))
	}
	g.curFn = fun

	bb := g.ctx.AddBasicBlock(g.curFn.fn, fmt.Sprintf("%s_bb", f.Var))
	g.builder.SetInsertPointAtEnd(bb)

	for i, b := range f.Bindings {
		g.genBindingAccesses(fun.fn.Param(i), b)
	}

	cpy := make(map[string]llvm.Value)
	for k, v := range fun.params {
		cpy[k] = v
	}

	for _, s := range f.Statements {
		g.generateStatement(cpy, s)
	}

}
