package lexer

import (
	"fmt"
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
	{
		input: `show inc(33)`,
		tokens: []Token{
			{Type: Show, Val: "show"},
			{Type: Variable, Val: "inc"},
			{Type: LParen, Val: "("},
			{Type: IntLiteral, Val: "33"},
			{Type: RParen, Val: ")"},
		},
	},
	{
		input: `show inc(3.3)`,
		tokens: []Token{
			{Type: Show, Val: "show"},
			{Type: Variable, Val: "inc"},
			{Type: LParen, Val: "("},
			{Type: FloatLiteral, Val: "3.3"},
			{Type: RParen, Val: ")"},
		},
	},
	{
		input: `show inc(3.)`,
		tokens: []Token{
			{Type: Show, Val: "show"},
			{Type: Variable, Val: "inc"},
			{Type: LParen, Val: "("},
			{Type: ILLEGAL, Val: "error: expected digits after the decimal received: )"},
		},
	},
	{
		input: `fn example(i : int, j : int) {
  return i / j
}`,
		tokens: []Token{
			{Type: Function, Val: "fn"},
			{Type: Variable, Val: "example"},
			{Type: LParen, Val: "("},
			{Type: Variable, Val: "i"},
			{Type: Colon, Val: ":"},
			{Type: Variable, Val: "int"},
			{Type: Comma, Val: ","},
			{Type: Variable, Val: "j"},
			{Type: Colon, Val: ":"},
			{Type: Variable, Val: "int"},
			{Type: RParen, Val: ")"},
			{Type: LCurly, Val: "{"},
			{Type: NewLine, Val: "\n"},
			{Type: Return, Val: "return"},
			{Type: Variable, Val: "i"},
			{Type: Divide, Val: "/"},
			{Type: Variable, Val: "j"},
			{Type: NewLine, Val: "\n"},
			{Type: RCurly, Val: "}"},
			{Type: EOF, Val: ""},
		},
	},
}

func TestLexer(t *testing.T) {
	for _, test := range test {
		l := NewLexer(test.input)

		i := 0
		for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
			fmt.Println(tok.DumpString())
			if i == len(test.tokens) {
				assert.Fail(t, "err, too many tokens, stopped on: %s\n\n NextToken: %+v", tok)
			}
			assert.Equal(t, test.tokens[i], tok)
			i++
		}
		fmt.Println("-------------")
	}
}
