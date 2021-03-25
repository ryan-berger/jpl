package main

import (
	"flag"
	"os"

	"github.com/ryan-berger/jpl/internal"
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
}

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	opts := []internal.CompilerOpts{internal.WithReader(file)}

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

	internal.NewCompiler(opts...).Compile()
}
