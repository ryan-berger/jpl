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
	fmt.Println(p.ParseProgram(false))
}

func TestParse_Fn(t *testing.T) {
	test := `
fn dense(input[W, H] : float[,], weights[Wi, Hi, Wo, Ho] : float[,,,]) : float[,] {
    assert W == Wi && H == Hi, "Weight matrix doesn't match input size"
    return array[i : W, j : H] \
      relu(sum[i2 : W, j2 : H] input[i2, j2] * weights[i2, j2, i, j])
}
`

	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	p.ParseProgram(false)
	assert.Empty(t, p.errors)
}

func TestParse_Cmd(t *testing.T) {
	test := `show x + 10`

	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	p.ParseProgram(false)
}
