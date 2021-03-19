package generator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/typed"
	"github.com/ryan-berger/jpl/internal/expander"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)

func expand(t *testing.T, program string) ast.Program {
	tokens, ok := lexer.Lex(program)
	assert.True(t, ok)

	tree, err := parser.Parse(tokens)
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	tree, err = typed.Check(tree)
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	tree, err = typed.Check(expander.Expand(tree))
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil
	}

	return tree
}

func TestData(t *testing.T) {
	program := expand(t, `let x = 10
print "hello world"`)
	p, mapper := dataSection(program)
	fmt.Println(p)
	textSection(program, mapper)
}
