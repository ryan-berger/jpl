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
		input: `}

test = 0

fn yeet`,
		tokens: []Token{
			{Type: RCurly, Val: "}"},
			{Type: NewLine, Val: "\n"},
			{Type: Variable, Val: "test"},
			{Type: Assign, Val: "="},
			{Type: IntLiteral, Val: "0"},
			{Type: NewLine, Val: "\n"},
			{Type: Function, Val: "fn"},
			{Type: Variable, Val: "yeet"},
			{Type: EOF},
		},
	},

	// TODO: actually make this one work....
	{
		input: `print "�"`,
		tokens: []Token{
			{Type: Print, Val: "print"},
			{Type: ILLEGAL, Val: "error, expected end quote received: ï"},
			{Type: EOF},
		},
	},
	{
		input: `"test string yeeet"`,
		tokens: []Token{
			{Type: String, Val: "\"test string yeeet\""},
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
		input: `show inc(3.3, \
3.3, 3.5)`,
		tokens: []Token{
			{Type: Show, Val: "show"},
			{Type: Variable, Val: "inc"},
			{Type: LParen, Val: "("},
			{Type: FloatLiteral, Val: "3.3"},
			{Type: Comma, Val: ","},
			{Type: FloatLiteral, Val: "3.3"},
			{Type: Comma, Val: ","},
			{Type: FloatLiteral, Val: "3.5"},
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
			{Type: EOF, Val: ""},
		},
	},
	{
		input: `// test comment
/*
test
test
test*/

fn example(i : int, j : int) {
  return i / j
}`,
		tokens: []Token{
			{Type: NewLine, Val: "\n"},
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
	{
		input: `-22
-1.3
+0.3`,
		tokens: []Token{
			{Type: IntLiteral, Val: "-22"},
			{Type: NewLine, Val: "\n"},
			{Type: FloatLiteral, Val: "-1.3"},
			{Type: NewLine, Val: "\n"},
			{Type: FloatLiteral, Val: "+0.3"},
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
				return
			}
			assert.Equal(t, test.tokens[i], tok)
			i++
		}
		fmt.Println("-------------")
	}
}
