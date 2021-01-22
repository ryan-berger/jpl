package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var test = []struct {
	input  string
	tokens []Token
}{
	{
		input: "{}[]*!",
		tokens: []Token{
			{Type: LCurly, Val: "{"},
			{Type: RCurly, Val: "}"},
			{Type: LBrace, Val: "["},
			{Type: RBrace, Val: "]"},
			{Type: Multiply, Val: "*"},
			{Type: Not, Val: "!"},
			{Type: EOF},
		},
	},
}

func TestLexer(t *testing.T) {
	for _, test := range test {
		l := NewLexer(test.input)
		for i := 0; i < len(test.tokens); i++ {
			assert.Equal(t,  test.tokens[i], l.NextToken())
		}
	}
}
