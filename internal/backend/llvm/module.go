package llvm

import (
	"fmt"
	"io"
	"os"

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

func Generate(p ast.Program, s *symbol.Table, w io.Writer) {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	ctx := llvm.NewContext()
	module := ctx.NewModule("main")
	builder := ctx.NewBuilder()

	defer builder.Dispose()

	g := generator{
		ctx:     ctx,
		builder: builder,
		module:  module,
		fns:     make(map[string]fn),
	}

	fpm := llvm.NewFunctionPassManagerForModule(module)
	fpm.AddCFGSimplificationPass()
	fpm.AddReassociatePass()
	fpm.AddInstructionCombiningPass()
	fpm.InitializeFunc()

	defer fpm.Dispose()

	g.genRuntime()
	g.generate(p, fpm)

	g.module.Dump()

	if err := llvm.VerifyModule(g.module, llvm.PrintMessageAction); err != nil {
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

func (g *generator) generate(p ast.Program, fpm llvm.PassManager) {
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
			//g.fns[command.Var].fn.Dump()
			//fpm.RunFunc(g.fns[command.Var].fn)
		default:
			cmds = append(cmds, cmd)
		}
	}

	fnType := llvm.FunctionType(g.ctx.Int64Type(), []llvm.Type{}, false)
	fn := llvm.AddFunction(g.module, "main", fnType)
	bb := g.ctx.AddBasicBlock(fn, "main_bb")
	g.builder.SetInsertPointAtEnd(bb)

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
		readStr := g.builder.CreateGlobalStringPtr(command.Src, "file_name")
		res := g.builder.CreateCall(g.fns["read_image"].fn, []llvm.Value{readStr}, "pict")

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
		fileName := g.builder.CreateGlobalStringPtr(command.Dest, "file_name")
		input := g.getExpr(vals, command.Expr)

		g.builder.CreateCall(g.fns["write_image"].fn, []llvm.Value{fileName, input}, "")
	case *ast.Show:
		typStr := command.Expr.Typ().String()
		str := g.builder.CreateGlobalStringPtr(typStr, "type")
		exp := g.getExpr(vals, command.Expr)

		ptr := g.builder.CreateAlloca(exp.Type(), "expr_ptr")
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
		args := make([]llvm.Value, len(fun.params))
		for i, e := range expr.Arguments {
			args[i] = g.getExpr(val, e)
		}
		return g.builder.CreateCall(fun.fn, args, "call")
	case *ast.SumTransform:
		return g.genSumTransform(val, expr)
	case *ast.ArrayTransform:
		return g.genArrayTransform(val, expr)
	case *ast.IfExpression:
		return g.genIf(val, expr)
	case *ast.TupleRefExpression:
		return g.builder.CreateExtractValue(g.getExpr(val, expr.Tuple), int(expr.Index), "tuple_lookup")
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
	}
	panic(fmt.Sprintf("unsupported expr: %T", expression))
}

func (g *generator) generateStatement(m map[string]llvm.Value, s ast.Statement) {
	switch stmt := s.(type) {
	case *ast.LetStatement:
		switch l := stmt.LValue.(type) {
		case *ast.Variable:
			exp := g.getExpr(m, stmt.Expr)
			m[l.Variable] = exp
		case *ast.VariableArr:
			exp := g.getExpr(m, stmt.Expr)
			m[l.Variable] = exp
			for i, v := range l.Variables {
				m[v] = g.builder.CreateExtractValue(exp, i, "rank")
			}
		}
	case *ast.ReturnStatement:
		g.builder.CreateRet(g.getExpr(m, stmt.Expr))
	case *ast.AssertStatement:
		assertFn := g.fns["fail_assertion"]
		exp := g.getExpr(m, stmt.Expr)

		failBB := g.ctx.AddBasicBlock(g.curFn.fn, "assertfail")
		contBB := g.ctx.AddBasicBlock(g.curFn.fn, "assertcont")

		g.builder.CreateCondBr(exp, failBB, contBB)

		str := g.builder.CreateGlobalStringPtr(stmt.Message, "assert")

		g.builder.SetInsertPointAtEnd(failBB)
		g.builder.CreateCall(assertFn.fn,
			[]llvm.Value{str},
			"",
		)
		g.builder.CreateUnreachable()
		g.builder.SetInsertPointAtEnd(contBB)
	default:
		panic(fmt.Sprintf("unsupported stmt: %T", s))
	}
}
