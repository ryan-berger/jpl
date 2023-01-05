package llvm

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

//go:embed "runtime/c/runtime.bc"
var runtime []byte

//go:embed "runtime/c/pngstuff.bc"
var pngstuff []byte

func getModule(ctx llvm.Context, name string, file []byte) llvm.Module {
	tf, err := os.CreateTemp("", fmt.Sprintf("%s.ll", name))
	if err != nil {
		panic(fmt.Sprintf("unable to create temp directory for runtime: %s", err))
	}

	tf.Write(file)
	tf.Close()

	path, err := filepath.Abs(tf.Name())
	if err != nil {
		panic(err)
	}

	mod, err := ctx.ParseBitcodeFile(path)
	if err != nil {
		panic(err)
	}

	os.Remove(path)

	return mod
}

func setupRuntime(ctx llvm.Context) llvm.Module {
	runtimeMod := getModule(ctx, "runtime", runtime)
	pngMod := getModule(ctx, "pngstuff", pngstuff)

	if err := llvm.LinkModules(runtimeMod, pngMod); err != nil {
		panic(err)
	}

	return runtimeMod
}

func Generate(p ast.Program, s *symbol.Table, w io.Writer) {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	ctx := llvm.NewContext()

	mod := setupRuntime(ctx)

	module := ctx.NewModule("main")
	builder := ctx.NewBuilder()

	defer builder.Dispose()

	g := generator{
		ctx:     ctx,
		builder: builder,
		module:  module,
		fns:     make(map[string]fn),
	}

	g.genRuntime()
	g.generate(p)

	passBuilder := llvm.NewPassManagerBuilder()

	passes := llvm.NewPassManager()
	defer passes.Dispose()
	passes.AddLICMPass()
	passes.AddGlobalDCEPass()
	passes.AddGlobalOptimizerPass()
	passes.AddIPSCCPPass()
	passes.AddAggressiveDCEPass()
	passes.AddFunctionAttrsPass()
	passes.AddFunctionInliningPass()
	passes.AddLoopUnrollPass()

	passBuilder.SetOptLevel(3)
	passBuilder.Populate(passes)

	passes.Run(module)

	if err := llvm.VerifyModule(module, llvm.PrintMessageAction); err != nil {
		panic(err)
	}

	if err := llvm.LinkModules(module, mod); err != nil {
		panic(err)
	}

	// output code

	trgt := llvm.DefaultTargetTriple()
	t, err := llvm.GetTargetFromTriple(trgt)
	if err != nil {
		panic(err)
	}

	tm := t.CreateTargetMachine(trgt, "generic", "",
		llvm.CodeGenLevelDefault, llvm.RelocDefault, llvm.CodeModelDefault)

	g.module.SetTarget(tm.Triple())
	g.module.SetDataLayout(tm.CreateTargetData().String())

	buf, err := tm.EmitToMemoryBuffer(g.module, llvm.ObjectFile)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("test.o", buf.Bytes(), 0666); err != nil {
		panic(err)
	}
}

func (g *generator) generate(p ast.Program) {
	for _, cmd := range p {
		if fn, ok := cmd.(*ast.Function); ok {
			g.declareFunction(fn)
		}
	}

	var cmds []ast.Command

	for _, cmd := range p {
		switch command := cmd.(type) {
		case *ast.Function:
			g.genFunction(command)
		default:
			cmds = append(cmds, cmd)
		}
	}

	fnType := llvm.FunctionType(g.ctx.Int64Type(), []llvm.Type{}, false)
	main := llvm.AddFunction(g.module, "main", fnType)
	bb := g.ctx.AddBasicBlock(main, "main_bb")
	g.builder.SetInsertPointAtEnd(bb)

	g.curFn = fn{
		fn: main,
	}

	globals := make(map[string]llvm.Value)
	for _, c := range cmds {
		g.generateCommand(globals, c)
	}
}

func (g *generator) generateCommand(vals map[string]llvm.Value, cmd ast.Command) {
	switch command := cmd.(type) {
	case *ast.Function:
		panic("function should have already been generated")
	case ast.Statement:
		g.generateStatement(vals, command)
	case *ast.Print:
		printStr := g.builder.CreateGlobalStringPtr(command.Str, "print")
		g.builder.CreateCall(g.fns["print"].fn, []llvm.Value{printStr}, "")
	case *ast.Read:
		readStr := g.getExpr(vals, command.Src)
		read := g.builder.CreateExtractValue(readStr, 0, "read")

		resPtr := g.builder.CreateMalloc(toLLVMType(g.ctx, types.Pict), "res_ptr")

		g.builder.CreateCall(g.fns["read_image"].fn, []llvm.Value{resPtr, read}, "")

		res := g.builder.CreateLoad(resPtr, "res")

		switch arg := command.Argument.(type) {
		case *ast.Variable:
			vals[arg.Variable] = res
		case *ast.VariableArr:
			vals[arg.Variable] = res
			for i, v := range arg.Variables {
				vals[v] = g.builder.CreateExtractValue(res, i, "rank")
			}
		}

	case *ast.Write:
		destStr := g.getExpr(vals, command.Dest)
		dest := g.builder.CreateExtractValue(destStr, 0, "dest_str")

		input := g.getExpr(vals, command.Expr)
		m := g.builder.CreateMalloc(input.Type(), "img_arg")
		g.builder.CreateStore(input, m)

		g.builder.CreateCall(g.fns["write_image"].fn, []llvm.Value{m, dest}, "")
	case *ast.Show:
		typStr := command.Expr.Typ().String()
		str := g.builder.CreateGlobalStringPtr(typStr, "type")
		exp := g.getExpr(vals, command.Expr)

		ptr := g.builder.CreateMalloc(exp.Type(), "expr_ptr")
		g.builder.CreateStore(exp, ptr)
		ptr = g.builder.CreateBitCast(ptr, llvm.PointerType(g.ctx.Int8Type(), 0), "cast")

		g.builder.CreateCall(g.fns["show"].fn, []llvm.Value{str, ptr}, "")

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
		return ctx.DoubleType()
	case p == types.Integer:
		return ctx.Int64Type()
	case p == types.Boolean:
		return ctx.Int1Type()
	case p == types.Str:
		return ctx.StructType([]llvm.Type{llvm.PointerType(ctx.Int8Type(), 0), ctx.Int64Type()}, false)
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

func (g *generator) genLVal(m map[string]llvm.Value, value ast.LValue, exp llvm.Value) {
	switch l := value.(type) {
	case *ast.LTuple:
		for i, lval := range l.Args {
			ev := g.builder.CreateExtractValue(exp, i, fmt.Sprintf("extract_tup_%d", i))
			g.genLVal(m, lval, ev)
		}
	case *ast.Variable:
		m[l.Variable] = exp
	case *ast.VariableArr:
		m[l.Variable] = exp
		for i, v := range l.Variables {
			m[v] = g.builder.CreateExtractValue(exp, i, "rank")
		}
	}
}

func (g *generator) generateStatement(m map[string]llvm.Value, s ast.Statement) {
	switch stmt := s.(type) {
	case *ast.LetStatement:
		exp := g.getExpr(m, stmt.Expr)
		g.genLVal(m, stmt.LValue, exp)
	case *ast.ReturnStatement:
		g.builder.CreateRet(g.getExpr(m, stmt.Expr))
	case *ast.AssertStatement:
		g.createAssert(g.getExpr(m, stmt.Expr), stmt.Message)
	default:
		panic(fmt.Sprintf("unsupported stmt: %T", s))
	}
}
