package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)

var debugLex bool
func init() {
	flag.BoolVar(&debugLex, "l", false, "lex")
	flag.BoolVar(&debugLex, "p", false, "parse")
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
	p.ParseProgram()
}
