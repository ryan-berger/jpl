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
fn gamma_decompress(x : float) : float {
   assert 0.0 <= x && x <= 1.0, "gamma_decompress argument out of range"
   return \
     if x <= 0.04045 then x / 12.92 else \
     pow((x + 0.055) / 1.055, 2.4)
}

read 
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
