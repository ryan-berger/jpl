package generator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/typed"
	"github.com/ryan-berger/jpl/internal/expander"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/ryan-berger/jpl/internal/symbol"
)

func expand(t *testing.T, program string) (ast.Program, *symbol.Table) {
	tokens, ok := lexer.Lex(program)
	assert.True(t, ok)

	tree, err := parser.Parse(tokens)
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil, nil
	}

	tree, _, err = typed.Check(tree)
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil, nil
	}

	tree, table, err := typed.Check(expander.Expand(tree))
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil, nil
	}

	return tree, table
}

func TestData(t *testing.T) {
	program, table := expand(t, `
read image "foo.png" to img

assert has_size(img, 10, 10), "Please do not modify foo.png, it should be 10x10"`)
	Generate(program, table, os.Stdout)
}