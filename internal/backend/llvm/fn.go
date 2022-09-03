package llvm

import (
	"fmt"
	"os"

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
		return llvm.StructType(collections.Map(bind.Bindings, curryBinding(ctx)), false)
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
		bind := f.(*ast.TypeBind)
		variable := bind.Argument.(*ast.Variable)

		llvmFn.Param(i).SetName(variable.Variable)
		m[variable.Variable] = llvmFn.Param(i)
	}

	g.fns[f.Var] = fn{
		fn:     llvmFn,
		params: m,
	}
}

func (g *generator) genFunction(f *ast.Function) {
	fun, ok := g.fns[f.Var]
	if !ok {
		panic(fmt.Sprintf("fn %s not found", f.Var))
	}
	g.curFn = fun
	fmt.Println(f.String())

	bb := g.ctx.AddBasicBlock(g.curFn.fn, fmt.Sprintf("%s_bb", f.Var))
	g.builder.SetInsertPointAtEnd(bb)

	cpy := make(map[string]llvm.Value)
	for k, v := range fun.params {
		cpy[k] = v
	}

	for _, s := range f.Statements {
		g.generateStatement(cpy, s)
	}

	if err := llvm.VerifyFunction(fun.fn, llvm.AbortProcessAction); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
