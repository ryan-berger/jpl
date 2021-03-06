package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/lexer"
)

func TestParse_Let(t *testing.T) {
	test := `let {x[W,H], y, z} = {(10*3)+10, 3, 4}`

	tokens, ok := lexer.Lex(test)
	assert.True(t, ok)


	fmt.Println(Parse(tokens))
}

func TestParse_Fn(t *testing.T) {
	test := `
print "test"
`
	tokens, ok := lexer.Lex(test)
	assert.True(t, ok)

	cmds, err := Parse(tokens)
	fmt.Println(cmds)
	assert.Nil(t, err)
}

func TestParse_Cmd(t *testing.T) {
	test := `show x + 10`

	tokens, ok := lexer.Lex(test)
	assert.True(t, ok)

	Parse(tokens)
}
