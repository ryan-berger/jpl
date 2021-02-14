package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)

var debugLex bool
var debugParse bool
var debug bool
func init() {
	flag.BoolVar(&debugLex, "l", false, "lex")
	flag.BoolVar(&debugParse, "p", false, "parse")
	flag.BoolVar(&debug, "d", false, "debug")
}

func main() {
	flag.Parse()

	file, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	l := lexer.NewLexer(string(file))
	tokens, success := l.LexAll()

	if debugLex {
		for _, tok := range tokens {
			fmt.Println(tok.DumpString())
		}

		if success {
			fmt.Println("Compilation succeeded")
		} else {
			fmt.Println("Compilation failed")
		}
		return
	}

	p := parser.NewParser(tokens)
	_, err = p.ParseProgram(debugParse)
	if err != nil {
		if debug {
			fmt.Println(err)
		}
		fmt.Println("Compilation failed")
		return
	}
	fmt.Println("Compilation succeeded")
}
