panic: identifier h not found

goroutine 1 [running]:
github.com/ryan-berger/jpl/internal/backend/llvm.(*generator).getExpr(0xc00010dbb8, 0xc0000061a0?, {0x54a7a8?, 0xc0000b30e0})
	/home/ryan/github.com/ryan-berger/jpl/internal/backend/llvm/module.go:190 +0x594
github.com/ryan-berger/jpl/internal/backend/llvm.(*generator).getExpr(0xc00010dbb8, 0x4db26d?, {0x54a828?, 0xc0000dccd0})
	/home/ryan/github.com/ryan-berger/jpl/internal/backend/llvm/module.go:169 +0x14b
github.com/ryan-berger/jpl/internal/backend/llvm.(*generator).generateStatement(0xc00010dbb8, 0x5?, {0x54a8a8?, 0xc0000b3a70?})
	/home/ryan/github.com/ryan-berger/jpl/internal/backend/llvm/module.go:201 +0xdb
github.com/ryan-berger/jpl/internal/backend/llvm.(*generator).genFunction(0xc00010dbb8, 0xc0000c2360)
	/home/ryan/github.com/ryan-berger/jpl/internal/backend/llvm/fn.go:114 +0x2f4
github.com/ryan-berger/jpl/internal/backend/llvm.(*generator).generate(0xc00010dbb8, {0xc0000f6280?, 0x7, 0xc000070bf8?}, {0x4c53c9?})
	/home/ryan/github.com/ryan-berger/jpl/internal/backend/llvm/module.go:66 +0x85
github.com/ryan-berger/jpl/internal/backend/llvm.Generate({0xc0000f6280, 0x7, 0x8}, 0x43cc25?, {0x629280?, 0x203000?})
	/home/ryan/github.com/ryan-berger/jpl/internal/backend/llvm/module.go:51 +0x24e
github.com/ryan-berger/jpl/internal.(*Compiler).generate(...)
	/home/ryan/github.com/ryan-berger/jpl/internal/compiler.go:141
github.com/ryan-berger/jpl/internal.(*Compiler).compile(0xc0000dc000)
	/home/ryan/github.com/ryan-berger/jpl/internal/compiler.go:175 +0x166
github.com/ryan-berger/jpl/internal.(*Compiler).Compile(0x507a00?)
	/home/ryan/github.com/ryan-berger/jpl/internal/compiler.go:180 +0x1d
main.main()
	/home/ryan/github.com/ryan-berger/jpl/cmd/compiler.go:105 +0x559
