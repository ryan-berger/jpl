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
		input: "&&||%{}[]*!/* */!=/* other comment */<=/**/>=/**/<>/**/==",
		tokens: []Token{
			{Type: And, Val: "&&", Line: 1},
			{Type: Or, Val: "||", Line: 1},
			{Type: Mod, Val: "%", Line: 1},
			{Type: LCurly, Val: "{", Line: 1},
			{Type: RCurly, Val: "}", Line: 1},
			{Type: LBrace, Val: "[", Line: 1},
			{Type: RBrace, Val: "]", Line: 1},
			{Type: Multiply, Val: "*", Line: 1},
			{Type: Not, Val: "!", Line: 1},
			{Type: NotEqualTo, Val: "!=", Line: 1},
			{Type: LessThanOrEqual, Val: "<=", Line: 1},
			{Type: GreaterThanOrEqual, Val: ">=", Line: 1},
			{Type: LessThan, Val: "<", Line: 1},
			{Type: GreaterThan, Val: ">", Line: 1},
			{Type: EqualTo, Val: "==", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `&|`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected '&' received '|'", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `|&`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected '|' received '&'", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `}

test = 0

fn yeet`,
		tokens: []Token{
			{Type: RCurly, Val: "}", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: Variable, Val: "test", Line: 3},
			{Type: Assign, Val: "=", Line: 3},
			{Type: IntLiteral, Val: "0", Line: 3},
			{Type: NewLine, Val: "\n", Line: 3},
			{Type: Function, Val: "fn", Line: 5},
			{Type: Variable, Val: "yeet", Line: 5},
			{Type: EOF},
		},
	},
	{
		input: `print "�"`,
		tokens: []Token{
			{Type: Print, Val: "print", Line: 1},
			{Type: ILLEGAL, Val: "error, expected end quote received: ï", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `"test string yeeet"`,
		tokens: []Token{
			{Type: String, Val: "\"test string yeeet\"", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `show inc(33)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1},
			{Type: Variable, Val: "inc", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: IntLiteral, Val: "33", Line: 1},
			{Type: RParen, Val: ")", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.3)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1},
			{Type: Variable, Val: "inc", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: FloatLiteral, Val: "3.3", Line: 1},
			{Type: RParen, Val: ")", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.3, \
3.3, 3.5)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1},
			{Type: Variable, Val: "inc", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: FloatLiteral, Val: "3.3", Line: 1},
			{Type: Comma, Val: ",", Line: 1},
			{Type: FloatLiteral, Val: "3.3", Line: 2},
			{Type: Comma, Val: ",", Line: 2},
			{Type: FloatLiteral, Val: "3.5", Line: 2},
			{Type: RParen, Val: ")", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1},
			{Type: Variable, Val: "inc", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: FloatLiteral, Val: "3.", Line: 1},
			{Type: RParen, Val: ")", Line: 1},
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
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: Function, Val: "fn", Line: 7},
			{Type: Variable, Val: "example", Line: 7},
			{Type: LParen, Val: "(", Line: 7},
			{Type: Variable, Val: "i", Line: 7},
			{Type: Colon, Val: ":", Line: 7},
			{Type: Int, Val: "int", Line: 7},
			{Type: Comma, Val: ",", Line: 7},
			{Type: Variable, Val: "j", Line: 7},
			{Type: Colon, Val: ":", Line: 7},
			{Type: Int, Val: "int", Line: 7},
			{Type: RParen, Val: ")", Line: 7},
			{Type: LCurly, Val: "{", Line: 7},
			{Type: NewLine, Val: "\n", Line: 7},
			{Type: Return, Val: "return", Line: 8},
			{Type: Variable, Val: "i", Line: 8},
			{Type: Divide, Val: "/", Line: 8},
			{Type: Variable, Val: "j", Line: 8},
			{Type: NewLine, Val: "\n", Line: 8},
			{Type: RCurly, Val: "}", Line: 9},
			{Type: EOF, Val: ""},
		},
	},
	{
		input: `fn example(i : int, j : int) {
  return i / j
}`,
		tokens: []Token{
			{Type: Function, Val: "fn", Line: 1},
			{Type: Variable, Val: "example", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: Variable, Val: "i", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Int, Val: "int", Line: 1},
			{Type: Comma, Val: ",", Line: 1},
			{Type: Variable, Val: "j", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Int, Val: "int", Line: 1},
			{Type: RParen, Val: ")", Line: 1},
			{Type: LCurly, Val: "{", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: Return, Val: "return", Line: 2},
			{Type: Variable, Val: "i", Line: 2},
			{Type: Divide, Val: "/", Line: 2},
			{Type: Variable, Val: "j", Line: 2},
			{Type: NewLine, Val: "\n", Line: 2},
			{Type: RCurly, Val: "}", Line: 3},
			{Type: EOF, Val: ""},
		},
	},
	{
		input: `\ `,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected newline after '\\' received: ' '", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `-22
-1.3
+0.3`,
		tokens: []Token{
			{Type: Minus, Val: "-", Line: 1},
			{Type: IntLiteral, Val: "22", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: Minus, Val: "-", Line: 2},
			{Type: FloatLiteral, Val: "1.3", Line: 2},
			{Type: NewLine, Val: "\n", Line: 2},
			{Type: Plus, Val: "+", Line: 3},
			{Type: FloatLiteral, Val: "0.3", Line: 3},
			{Type: EOF},
		},
	},
	{
		input: `print ""`,
		tokens: []Token{
			{Type: Print, Val: "print", Line: 1},
			{Type: String, Val: "\"\"", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: "1234567\x0012345",
		tokens: []Token{
			{Type: IntLiteral, Val: "1234567", Line: 1},
			{Type: ILLEGAL, Val: "error, received invalid character: \x00", Line: 1},
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
		input: "/* \x00*/",
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected closing '*/' received invalid character: '\x00'", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `/*`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected closing '*/' received EOF", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: "\x00",
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, received invalid character: \x00", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: "print \"\x00\"",
		tokens: []Token{
			{Type: Print, Val: "print", Line: 1},
			{Type: ILLEGAL, Val: "error, expected end quote received: \x00", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `fn sphere_point({x : float, y : float}) : float3 {
    // p(t) = { 4 - t, t * x / 4, t * y / 4 }
    let r =`,
		tokens: []Token{
			{Type: Function, Val: "fn", Line: 1},
			{Type: Variable, Val: "sphere_point", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: LCurly, Val: "{", Line: 1},
			{Type: Variable, Val: "x", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Float, Val: "float", Line: 1},
			{Type: Comma, Val: ",", Line: 1},
			{Type: Variable, Val: "y", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Float, Val: "float", Line: 1},
			{Type: RCurly, Val: "}", Line: 1},
			{Type: RParen, Val: ")", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Float3, Val: "float3", Line: 1},
			{Type: LCurly, Val: "{", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: Let, Val: "let", Line: 3},
			{Type: Variable, Val: "r", Line: 3},
			{Type: Assign, Val: "=", Line: 3},
			{Type: EOF},
		},
	},
	{
		input: `let r = 10 // test comment`,
		tokens: []Token{
			{Type: Let, Val: "let", Line: 1},
			{Type: Variable, Val: "r", Line: 1},
			{Type: Assign, Val: "=", Line: 1},
			{Type: IntLiteral, Val: "10", Line: 1},
			{Type: EOF},
		},
	},
	{
		input: `M] ( \
      /* by using clamp here, we are "extending" the boundary pixels of a */ \
      k`,
		tokens: []Token{
			{Type: Variable, Val: "M", Line: 1},
			{Type: RBrace, Val: "]", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: Variable, Val: "k", Line: 3},
			{Type: EOF},
		},
	},
	{
		input: `+ mi /* Green is a good lazy monochrome */
}`,
		tokens: []Token{
			{Type: Plus, Val: "+", Line: 1},
			{Type: Variable, Val: "mi", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: RCurly, Val: "}", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `+ mi // Green is a good lazy monochrome
}`,
		tokens: []Token{
			{Type: Plus, Val: "+", Line: 1},
			{Type: Variable, Val: "mi", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: RCurly, Val: "}", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `.31111 // Green is a good lazy monochrome
}`,
		tokens: []Token{
			{Type: FloatLiteral, Val: ".31111", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: RCurly, Val: "}", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `fn my_fn(x : int) : int { return x }`,
		tokens: []Token{
			{Type: Function, Val: "fn", Line: 1},
			{Type: Variable, Val: "my_fn", Line: 1},
			{Type: LParen, Val: "(", Line: 1},
			{Type: Variable, Val: "x", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Int, Val: "int", Line: 1},
			{Type: RParen, Val: ")", Line: 1},
			{Type: Colon, Val: ":", Line: 1},
			{Type: Int, Val: "int", Line: 1},
			{Type: LCurly, Val: "{", Line: 1},
			{Type: Return, Val: "return", Line: 1},
			{Type: Variable, Val: "x", Line: 1},
			{Type: RCurly, Val: "}", Line: 1},
			{Type: EOF},
		},

	},
	{
		input: "return .38\ntime write",
		tokens: []Token{
			{Type: Return, Val: "return", Line: 1},
			{Type: FloatLiteral, Val: ".38", Line: 1},
			{Type: NewLine, Val: "\n", Line: 1},
			{Type: Time, Val: "time", Line: 2},
			{Type: Write, Val: "write", Line: 2},
			{Type: EOF},
		},
	},
}

func TestLexer(t *testing.T) {
	for i, test := range tests {
		fmt.Println(i)
		l := NewLexer(test.input)
		tokens, _ := l.LexAll()
		for _, tok := range tokens {
			fmt.Println(tok.DumpString())
		}
		assert.Equal(t, test.tokens, tokens)
		fmt.Println("-------------")
	}
}
