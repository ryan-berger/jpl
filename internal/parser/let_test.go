package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/lexer"
)

func TestParse_Let(t *testing.T) {
	test := `let {x[W,H], y, z} = {(10*3)+10, 3, 4}`

	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	assert.Empty(t, p.errors)
	fmt.Println(p.ParseProgram())
}

func TestParse_Fn(t *testing.T) {
	test := `
fn my_fn(x : int) : int[,][][] {
  assert 1 == 1, "my str"
  let x = array[i : 10] 10
  let y = myfn(1, 2, 3)
  return x 
}`

	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	p.ParseProgram()
	assert.Empty(t, p.errors)
}

func TestParse_Cmd(t *testing.T) {
	test := `show x + 10`

	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	p.ParseProgram()
}
