package nasm

import (
	"os"
	"testing"

	"github.com/ryan-berger/jpl/internal/ast/flatten"
	"github.com/ryan-berger/jpl/internal/ast/types/checker"
	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/ryan-berger/jpl/internal/symbol"
)

func parseAndFlatten(t *testing.T, program string) (ast.Program, *symbol.Table) {
	tokens, ok := lexer.Lex(program)
	assert.True(t, ok)

	tree, err := parser.Parse(tokens)
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil, nil
	}

	tree, _, err = checker.Check(tree)
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil, nil
	}


	tree, table, err := checker.Check(flatten.Flatten(tree))
	if err != nil {
		assert.FailNow(t, err.Error())
		return nil, nil
	}

	return tree, table
}

func TestData(t *testing.T) {
	program, table := parseAndFlatten(t, `
read image "foo.png" to img
write image blur(img, 3.14) to "foo.png"

let x = if 1 == 1 then 0 else 2
fn test() : {} {
  let y = 10
  let z = 22
  return {}
}

let g = test()
`)


	Generate(program, table, os.Stdout)
}
