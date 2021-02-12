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
	assert.Nil(t, p.error)
	fmt.Println(p.ParseProgram(false))
}

func TestParse_Fn(t *testing.T) {
	test := `
write image C() to "E" 
print "8" 
assert 62,
write image 6 to "s*2x[i"
show .3694592063
assert ""
show 09
time read image "" i
read "" q2V. to
return 44
`
	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	cmds, err := p.ParseProgram(false)
	fmt.Println(cmds)
	assert.Nil(t, err)
}

func TestParse_Cmd(t *testing.T) {
	test := `show x + 10`

	tokens, ok := lexer.NewLexer(test).LexAll()
	assert.True(t, ok)

	p := NewParser(tokens)
	p.ParseProgram(false)
}
