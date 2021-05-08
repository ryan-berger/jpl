package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/typed"
	"github.com/ryan-berger/jpl/internal/expander"
	"github.com/ryan-berger/jpl/internal/generator"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/optimizer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/ryan-berger/jpl/internal/symbol"
)

type PrintMode int

const (
	_ PrintMode = iota
	Lex
	Parse
	TypeCheck
	Flatten
	ASM
)

type Compiler struct {
	input         io.Reader
	output        io.Writer
	mode          PrintMode
	optimizations []optimizer.Optimization
}

type CompilerOpts func(c *Compiler)

func WithPrintMode(p PrintMode) CompilerOpts {
	return func(c *Compiler) {
		c.mode = p
	}
}

func WithWriter(w io.Writer) CompilerOpts {
	return func(c *Compiler) {
		c.output = w
	}
}

func WithReader(r io.Reader) CompilerOpts {
	return func(c *Compiler) {
		c.input = r
	}
}

func WithOptimizations(optimizations []optimizer.Optimization) CompilerOpts {
	return func(c *Compiler) {
		c.optimizations = optimizations
	}
}

func NewCompiler(opts ...CompilerOpts) *Compiler {
	c := &Compiler{
		input:  os.Stdin,
		output: ioutil.Discard,
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

func (c *Compiler) typeCheck(program ast.Program) (ast.Program, error) {
	newProgram, _, err := typed.Check(program)
	if err != nil {
		return nil, err
	}

	if c.mode == TypeCheck {
		fmt.Println(newProgram.SExpr())
	}

	return newProgram, nil
}

func (c *Compiler) expand(program ast.Program) (ast.Program, *symbol.Table) {
	expanded := expander.Expand(program)

	_, table, err := typed.Check(expanded)
	if err != nil {
		panic(fmt.Sprintf("nice, you really messed up, %s", err))
	}

	return expanded, table
}

func (c *Compiler) generate(program ast.Program, table *symbol.Table) {
	if c.mode == ASM {
		c.output = os.Stdout
	}
	generator.Generate(program, table, c.output)
}

func (c *Compiler) compile() error {
	tokens, err := c.lex()
	if err != nil || c.mode == Lex {
		return err
	}

	program, err := c.parse(tokens)
	if err != nil || c.mode == Parse {
		return err
	}

	program, err = c.typeCheck(program)
	if err != nil || c.mode == TypeCheck {
		return err
	}

	if len(c.optimizations) != 0 {
		for _, o := range c.optimizations {
			program = o(program)
		}
		fmt.Println(program.SExpr())
		return nil
	}

	program, table := c.expand(program)
	if c.mode == Flatten {
		fmt.Println(program.SExpr())
		return nil
	}

	c.generate(program, table)
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
