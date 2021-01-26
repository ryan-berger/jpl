package lexer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	input  string
	tokens []Token
}{
	{
		input: "&&||%{}[]*!/* */!=/* other comment */<=/**/>=/**/<>",
		tokens: []Token{
			{Type: And, Val: "&&"},
			{Type: Or, Val: "||"},
			{Type: Mod, Val: "%"},
			{Type: LCurly, Val: "{"},
			{Type: RCurly, Val: "}"},
			{Type: LBrace, Val: "["},
			{Type: RBrace, Val: "]"},
			{Type: Multiply, Val: "*"},
			{Type: Not, Val: "!"},
			{Type: NotEqualTo, Val: "!="},
			{Type: LessThanOrEqual, Val: "<="},
			{Type: GreaterThanOrEqual, Val: ">="},
			{Type: LessThan, Val: "<"},
			{Type: GreaterThan, Val: ">"},
			{Type: EOF},
		},
	},
	{
		input: `&|`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected '&' received '|'"},
			{Type: EOF},
		},
	},
	{
		input: `|&`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected '|' received '&'"},
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
			{Type: EOF},
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
			{Type: EOF},
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
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.)`,
		tokens: []Token{
			{Type: Show, Val: "show"},
			{Type: Variable, Val: "inc"},
			{Type: LParen, Val: "("},
			{Type: FloatLiteral, Val: "3."},
			{Type: RParen, Val: ")"},
			{Type: EOF},
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
		input: `\ `,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected newline after '\\' received: ' '"},
			{Type: EOF},
		},
	},
	{
		input: `-22
-1.3
+0.3`,
		tokens: []Token{
			{Type: Minus, Val: "-"},
			{Type: IntLiteral, Val: "22"},
			{Type: NewLine, Val: "\n"},
			{Type: Minus, Val: "-"},
			{Type: FloatLiteral, Val: "1.3"},
			{Type: NewLine, Val: "\n"},
			{Type: Plus, Val: "+"},
			{Type: FloatLiteral, Val: "0.3"},
			{Type: EOF},
		},
	},
	{
		input: `print ""`,
		tokens: []Token{
			{Type: Print, Val: "print"},
			{Type: String, Val: "\"\""},
			{Type: EOF},
		},
	},
	{
		input: `//000000000000`,
		tokens: []Token{
			{Type: EOF},
		},
	},
	{
		input: `/* */`,
		tokens: []Token{
			{Type: EOF},
		},
	},
	{
		input: `/*`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected closing '*/' received EOF"},
			{Type: EOF},
		},
	},
	{
		input: "\x00",
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, received invalid character: \x00"},
			{Type: EOF},
		},
	},
	{
		input: "print \"\x00\"",
		tokens: []Token{
			{Type: Print, Val: "print"},
			{Type: ILLEGAL, Val: "error, expected end quote received: \x00"},
			{Type: EOF},
		},
	},
	{
		input: `fn sphere_point({x : float, y : float}) : float3 {
    // p(t) = { 4 - t, t * x / 4, t * y / 4 }
    let r =`,
		tokens: []Token{
			{Type: Function, Val: "fn"},
			{Type: Variable, Val: "sphere_point"},
			{Type: LParen, Val: "("},
			{Type: LCurly, Val: "{"},
			{Type: Variable, Val: "x"},
			{Type: Colon, Val: ":"},
			{Type: Variable, Val: "float"},
			{Type: Comma, Val: ","},
			{Type: Variable, Val: "y"},
			{Type: Colon, Val: ":"},
			{Type: Variable, Val: "float"},
			{Type: RCurly, Val: "}"},
			{Type: RParen, Val: ")"},
			{Type: Colon, Val: ":"},
			{Type: Variable, Val: "float3"},
			{Type: LCurly, Val: "{"},
			{Type: NewLine, Val: "\n"},
			{Type: Let, Val: "let"},
			{Type: Variable, Val: "r"},
			{Type: Assign, Val: "="},
			{Type: EOF},
		},
	},
	{
		input: `let r = 10 // test comment`,
		tokens: []Token{
			{Type: Let, Val: "let"},
			{Type: Variable, Val: "r"},
			{Type: Assign, Val: "="},
			{Type: IntLiteral, Val: "10"},
			{Type: EOF},
		},
	},
	{
		input: `M] ( \
      /* by using clamp here, we are "extending" the boundary pixels of a */ \
      k`,
		tokens: []Token{
			{Type: Variable, Val: "M"},
			{Type: RBrace, Val: "]"},
			{Type: LParen, Val: "("},
			{Type: Variable, Val: "k"},
			{Type: EOF},
		},
	},
}

func TestLexer(t *testing.T) {
	for _, test := range tests {
		l := NewLexer(test.input)
		tokens, _ := l.LexAll()
		for _, tok := range tokens {
			fmt.Println(tok.DumpString())
		}
		assert.Equal(t, test.tokens, tokens)
		fmt.Println("-------------")
	}
}
