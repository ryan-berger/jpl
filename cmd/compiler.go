package main

import (
	"flag"
	"os"

	"github.com/ryan-berger/jpl/internal"
	"github.com/ryan-berger/jpl/internal/optimizer"
)

var debugLex bool
var debugParse bool
var debug bool
var typed bool
var flattened bool
var asm bool

var outFile string

func init() {
	flag.BoolVar(&debugLex, "l", false, "lex")
	flag.BoolVar(&debugParse, "p", false, "parse")
	flag.BoolVar(&debug, "d", false, "debug")
	flag.BoolVar(&typed, "t", false, "types")
	flag.BoolVar(&flattened, "f", false, "flatten")
	flag.BoolVar(&asm, "s", false, "flatten")
	flag.StringVar(&outFile, "o", "", "out file")
	flag.Bool("cf", false, "cf")
	flag.Bool("cp", false, "cp")
	flag.Bool("dce", false, "dce")
	flag.Bool("peep", false, "peep")
}

var optimization = map[string]optimizer.Optimization{
	"-cf": optimizer.ConstantFold,
	"-cp": optimizer.ConstantProp,
	"-dce": optimizer.DeadCode,
	"-peep": optimizer.Peephole,
}

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	var optimizations []optimizer.Optimization
	for _, f := range os.Args {
		if o, ok := optimization[f]; ok {
			optimizations = append(optimizations, o)
		}
	}

	opts := []internal.CompilerOpts{
		internal.WithReader(file),
		internal.WithOptimizations(optimizations),
	}

	var mode internal.PrintMode
	switch {
	case debugLex:
		mode = internal.Lex
	case debugParse:
		mode = internal.Parse
	case typed:
		mode = internal.TypeCheck
	case flattened:
		mode = internal.Flatten
	case asm:
		mode = internal.ASM
	}

	if mode != 0 {
		opts = append(opts, internal.WithPrintMode(mode))
	}

	if outFile != "" {
		file, err := os.Create(outFile)
		if err != nil {
			panic(err)
		}
		file.Truncate(0)
		file.Seek(0, 0)

		defer file.Close()
		opts = append(opts, internal.WithWriter(file))
	}

	internal.
		NewCompiler(opts...).
		Compile()
}
