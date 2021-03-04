package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/typed"
	"github.com/ryan-berger/jpl/internal/expander"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)

type PrintMode int

const (
	_ PrintMode = iota
	Lex
	Parse
	TypeCheck
)

type Compiler struct {
	input  io.Reader
	output io.Writer
	mode   PrintMode
}

type CompilerOpts func(c *Compiler)

func WithPrintMode(p PrintMode) CompilerOpts {
	return func(c *Compiler) {
		c.mode = p
	}
}

func WithReader(r io.Reader) CompilerOpts {
	return func(c *Compiler) {
		c.input = r
	}
}

func NewCompiler(opts ...CompilerOpts) *Compiler {
	c := &Compiler{
		input:  os.Stdin,
		output: os.Stdout,
	}

	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Compiler) lex() ([]lexer.Token, error) {
	b, err := ioutil.ReadAll(c.input)
	if err != nil {
		return nil, err
	}

	tokens, ok := lexer.Lex(string(b))
	if c.mode == Lex {
		for _, t := range tokens {
			fmt.Println(t.DumpString())
		}
	}
	if !ok { // TODO: print out real errors
		return nil, fmt.Errorf("")
	}

	return tokens, nil
}

func (c *Compiler) parse(tokens []lexer.Token) (ast.Program, error) {
	program, err := parser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	if c.mode == Parse {
		fmt.Println(program.SExpr())
	}
	return program, nil
}

func (c *Compiler) compile() error {
	tokens, err := c.lex()
	if err != nil {
		return err
	}

	if c.mode == Lex {
		return nil
	}

	program, err := c.parse(tokens)
	if err != nil {
		return err
	}

	if c.mode == Parse {
		return nil
	}

	program, err = typed.Check(program)

	if err != nil {
		return err
	}

	expanded := expander.Expand(program)

	_, err = typed.Check(expanded)
	if err != nil {
		panic("nice, you really messed up")
	}

	return nil
}

func (c *Compiler) Compile() {
	if err := c.compile(); err != nil {
		fmt.Println(err)
		fmt.Println("Compilation failed")
		return
	}

	fmt.Println("Compilation succeeded")
}
