package lexer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
			{Type: And, Val: "&&", Line: 1, Character: 0},
			{Type: Or, Val: "||", Line: 1, Character: 2},
			{Type: Mod, Val: "%", Line: 1, Character: 4},
			{Type: LCurly, Val: "{", Line: 1, Character: 5},
			{Type: RCurly, Val: "}", Line: 1, Character: 6},
			{Type: LBrace, Val: "[", Line: 1, Character: 7},
			{Type: RBrace, Val: "]", Line: 1, Character: 8},
			{Type: Multiply, Val: "*", Line: 1, Character: 9},
			{Type: Not, Val: "!", Line: 1, Character: 10},
			{Type: NotEqualTo, Val: "!=", Line: 1, Character: 16},
			{Type: LessThanOrEqual, Val: "<=", Line: 1, Character: 37},
			{Type: GreaterThanOrEqual, Val: ">=", Line: 1, Character: 43},
			{Type: LessThan, Val: "<", Line: 1, Character: 49},
			{Type: GreaterThan, Val: ">", Line: 1, Character: 50},
			{Type: EqualTo, Val: "==", Line: 1, Character: 55},
			{Type: EOF},
		},
	},
	{
		input: `&|`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected '&' received '|' at line 1 position 1", Line: 1, Character: 1},
			{Type: EOF},
		},
	},
	{
		input: `|&`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected '|' received '&' at line 1 position 1", Line: 1, Character: 1},
			{Type: EOF},
		},
	},
	{
		input: `}

test = 0

fn yeet`,
		tokens: []Token{
			{Type: RCurly, Val: "}", Line: 1, Character: 0},
			{Type: NewLine, Val: "\n", Line: 1, Character: 1},
			{Type: Variable, Val: "test", Line: 3, Character: 0},
			{Type: Assign, Val: "=", Line: 3, Character: 5},
			{Type: IntLiteral, Val: "0", Line: 3, Character: 7},
			{Type: NewLine, Val: "\n", Line: 3, Character: 8},
			{Type: Function, Val: "fn", Line: 5, Character: 0},
			{Type: Variable, Val: "yeet", Line: 5, Character: 3},
			{Type: EOF},
		},
	},
	{
		input: `print "�"`,
		tokens: []Token{
			{Type: Print, Val: "print", Line: 1, Character: 0},
			{Type: ILLEGAL, Val: "error, expected end quote received: ï at line 1 position 8", Line: 1, Character: 8},
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
			{Type: Show, Val: "show", Line: 1, Character: 0},
			{Type: Variable, Val: "inc", Line: 1, Character: 5},
			{Type: LParen, Val: "(", Line: 1, Character: 8},
			{Type: IntLiteral, Val: "33", Line: 1, Character: 9},
			{Type: RParen, Val: ")", Line: 1, Character: 11},
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.3)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1, Character: 0},
			{Type: Variable, Val: "inc", Line: 1, Character: 5},
			{Type: LParen, Val: "(", Line: 1, Character: 8},
			{Type: FloatLiteral, Val: "3.3", Line: 1, Character: 9},
			{Type: RParen, Val: ")", Line: 1, Character: 12},
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.3, \
3.3, 3.5)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1, Character: 0},
			{Type: Variable, Val: "inc", Line: 1, Character: 5},
			{Type: LParen, Val: "(", Line: 1, Character: 8},
			{Type: FloatLiteral, Val: "3.3", Line: 1, Character: 9},
			{Type: Comma, Val: ",", Line: 1, Character: 12},
			{Type: FloatLiteral, Val: "3.3", Line: 2, Character: 0},
			{Type: Comma, Val: ",", Line: 2, Character: 3},
			{Type: FloatLiteral, Val: "3.5", Line: 2, Character: 5},
			{Type: RParen, Val: ")", Line: 2, Character: 8},
			{Type: EOF},
		},
	},
	{
		input: `show inc(3.)`,
		tokens: []Token{
			{Type: Show, Val: "show", Line: 1, Character: 0},
			{Type: Variable, Val: "inc", Line: 1, Character: 5},
			{Type: LParen, Val: "(", Line: 1, Character: 8},
			{Type: FloatLiteral, Val: "3.", Line: 1, Character: 9},
			{Type: RParen, Val: ")", Line: 1, Character: 11},
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
			{Type: NewLine, Val: "\n", Line: 1, Character: 15},
			{Type: Function, Val: "fn", Line: 7, Character: 0},
			{Type: Variable, Val: "example", Line: 7, Character: 3},
			{Type: LParen, Val: "(", Line: 7, Character: 10},
			{Type: Variable, Val: "i", Line: 7, Character: 11},
			{Type: Colon, Val: ":", Line: 7, Character: 13},
			{Type: Int, Val: "int", Line: 7, Character: 15},
			{Type: Comma, Val: ",", Line: 7, Character: 18},
			{Type: Variable, Val: "j", Line: 7, Character: 20},
			{Type: Colon, Val: ":", Line: 7, Character: 22},
			{Type: Int, Val: "int", Line: 7, Character: 24},
			{Type: RParen, Val: ")", Line: 7, Character: 27},
			{Type: LCurly, Val: "{", Line: 7, Character: 29},
			{Type: NewLine, Val: "\n", Line: 7, Character: 30},
			{Type: Return, Val: "return", Line: 8, Character: 2},
			{Type: Variable, Val: "i", Line: 8, Character: 9},
			{Type: Divide, Val: "/", Line: 8, Character: 11},
			{Type: Variable, Val: "j", Line: 8, Character: 13},
			{Type: NewLine, Val: "\n", Line: 8, Character: 14},
			{Type: RCurly, Val: "}", Line: 9, Character: 0},
			{Type: EOF, Val: "", Line: 0, Character: 0},
		},
	},
	{
		input: `fn example(i : int, j : int) {
  return i / j
}`,
		tokens: []Token{
			{Type: Function, Val: "fn", Line: 1, Character: 0},
			{Type: Variable, Val: "example", Line: 1, Character: 3},
			{Type: LParen, Val: "(", Line: 1, Character: 10},
			{Type: Variable, Val: "i", Line: 1, Character: 11},
			{Type: Colon, Val: ":", Line: 1, Character: 13},
			{Type: Int, Val: "int", Line: 1, Character: 15},
			{Type: Comma, Val: ",", Line: 1, Character: 18},
			{Type: Variable, Val: "j", Line: 1, Character: 20},
			{Type: Colon, Val: ":", Line: 1, Character: 22},
			{Type: Int, Val: "int", Line: 1, Character: 24},
			{Type: RParen, Val: ")", Line: 1, Character: 27},
			{Type: LCurly, Val: "{", Line: 1, Character: 29},
			{Type: NewLine, Val: "\n", Line: 1, Character: 30},
			{Type: Return, Val: "return", Line: 2, Character: 2},
			{Type: Variable, Val: "i", Line: 2, Character: 9},
			{Type: Divide, Val: "/", Line: 2, Character: 11},
			{Type: Variable, Val: "j", Line: 2, Character: 13},
			{Type: NewLine, Val: "\n", Line: 2, Character: 14},
			{Type: RCurly, Val: "}", Line: 3, Character: 0},
			{Type: EOF, Val: "", Line: 0, Character: 0},
		},
	},
	{
		input: `\ `,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected newline after '\\' received: ' ' at line 1 position 1", Line: 1, Character: 1},
			{Type: EOF},
		},
	},
	{
		input: `-22
-1.3
+0.3`,
		tokens: []Token{
			{Type: Minus, Val: "-", Line: 1, Character: 0},
			{Type: IntLiteral, Val: "22", Line: 1, Character: 1},
			{Type: NewLine, Val: "\n", Line: 1, Character: 3},
			{Type: Minus, Val: "-", Line: 2, Character: 0},
			{Type: FloatLiteral, Val: "1.3", Line: 2, Character: 1},
			{Type: NewLine, Val: "\n", Line: 2, Character: 4},
			{Type: Plus, Val: "+", Line: 3, Character: 0},
			{Type: FloatLiteral, Val: "0.3", Line: 3, Character: 1},
			{Type: EOF},
		},
	},
	{
		input: `print ""`,
		tokens: []Token{
			{Type: Print, Val: "print", Line: 1, Character: 0},
			{Type: String, Val: "\"\"", Line: 1, Character: 6},
			{Type: EOF},
		},
	},
	{
		input: "1234567\x0012345",
		tokens: []Token{
			{Type: IntLiteral, Val: "1234567", Line: 1, Character: 0},
			{Type: ILLEGAL, Val: "error, received invalid character: \x00 at line 1 position 8", Line: 1, Character: 8},
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
			{Type: ILLEGAL, Val: "error, expected closing '*/' received invalid character: '\x00' at line 1 position 4", Line: 1, Character: 4},
			{Type: EOF},
		},
	},
	{
		input: `/*`,
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, expected closing '*/' received EOF at line 1 position 3", Line: 1, Character: 3},
			{Type: EOF},
		},
	},
	{
		input: "\x00",
		tokens: []Token{
			{Type: ILLEGAL, Val: "error, received invalid character: \x00 at line 1 position 1", Line: 1, Character: 1},
			{Type: EOF},
		},
	},
	{
		input: "print \"\x00\"",
		tokens: []Token{
			{Type: Print, Val: "print", Line: 1},
			{Type: ILLEGAL, Val: "error, expected end quote received: \x00 at line 1 position 8", Line: 1, Character: 8},
			{Type: EOF},
		},
	},
	{
		input: `fn sphere_point({x : float, y : float}) : float3 {
    // p(t) = { 4 - t, t * x / 4, t * y / 4 }
    let r =`,
		tokens: []Token{
			{Type: Function, Val: "fn", Line: 1, Character: 0},
			{Type: Variable, Val: "sphere_point", Line: 1, Character: 3},
			{Type: LParen, Val: "(", Line: 1, Character: 15},
			{Type: LCurly, Val: "{", Line: 1, Character: 16},
			{Type: Variable, Val: "x", Line: 1, Character: 17},
			{Type: Colon, Val: ":", Line: 1, Character: 19},
			{Type: Float, Val: "float", Line: 1, Character: 21},
			{Type: Comma, Val: ",", Line: 1, Character: 26},
			{Type: Variable, Val: "y", Line: 1, Character: 28},
			{Type: Colon, Val: ":", Line: 1, Character: 30},
			{Type: Float, Val: "float", Line: 1, Character: 32},
			{Type: RCurly, Val: "}", Line: 1, Character: 37},
			{Type: RParen, Val: ")", Line: 1, Character: 38},
			{Type: Colon, Val: ":", Line: 1, Character: 40},
			{Type: Float3, Val: "float3", Line: 1, Character: 42},
			{Type: LCurly, Val: "{", Line: 1, Character: 49},
			{Type: NewLine, Val: "\n", Line: 1, Character: 50},
			{Type: Let, Val: "let", Line: 3, Character: 4},
			{Type: Variable, Val: "r", Line: 3, Character: 8},
			{Type: Assign, Val: "=", Line: 3, Character: 10},
			{Type: EOF, Val: "", Line: 0, Character: 0},
		},
	},
	{
		input: `let r = 10 // test comment`,
		tokens: []Token{
			{Type: Let, Val: "let", Line: 1, Character: 0},
			{Type: Variable, Val: "r", Line: 1, Character: 4},
			{Type: Assign, Val: "=", Line: 1, Character: 6},
			{Type: IntLiteral, Val: "10", Line: 1, Character: 8},
			{Type: EOF},
		},
	},
	{
		input: `M] ( \
      /* by using clamp here, we are "extending" the boundary pixels of a */ \
      k`,
		tokens: []Token{
			{Type: Variable, Val: "M", Line: 1},
			{Type: RBrace, Val: "]", Line: 1, Character: 1},
			{Type: LParen, Val: "(", Line: 1, Character: 3},
			{Type: Variable, Val: "k", Line: 3, Character: 6},
			{Type: EOF},
		},
	},
	{
		input: `+ mi /* Green is a good lazy monochrome */
}`,
		tokens: []Token{
			{Type: Plus, Val: "+", Line: 1, Character: 0},
			{Type: Variable, Val: "mi", Line: 1, Character: 2},
			{Type: NewLine, Val: "\n", Line: 1, Character: 42},
			{Type: RCurly, Val: "}", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `+ mi // Green is a good lazy monochrome
}`,
		tokens: []Token{
			{Type: Plus, Val: "+", Line: 1},
			{Type: Variable, Val: "mi", Line: 1, Character: 2},
			{Type: NewLine, Val: "\n", Line: 1, Character: 39},
			{Type: RCurly, Val: "}", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `.31111 // Green is a good lazy monochrome
}`,
		tokens: []Token{
			{Type: FloatLiteral, Val: ".31111", Line: 1, Character: 0},
			{Type: NewLine, Val: "\n", Line: 1, Character: 41},
			{Type: RCurly, Val: "}", Line: 2},
			{Type: EOF},
		},
	},
	{
		input: `fn my_fn(x : int) : int { return x }`,
		tokens: []Token{
			{Type: Function, Val: "fn", Line: 1, Character: 0},
			{Type: Variable, Val: "my_fn", Line: 1, Character: 3},
			{Type: LParen, Val: "(", Line: 1, Character: 8},
			{Type: Variable, Val: "x", Line: 1, Character: 9},
			{Type: Colon, Val: ":", Line: 1, Character: 11},
			{Type: Int, Val: "int", Line: 1, Character: 13},
			{Type: RParen, Val: ")", Line: 1, Character: 16},
			{Type: Colon, Val: ":", Line: 1, Character: 18},
			{Type: Int, Val: "int", Line: 1, Character: 20},
			{Type: LCurly, Val: "{", Line: 1, Character: 24},
			{Type: Return, Val: "return", Line: 1, Character: 26},
			{Type: Variable, Val: "x", Line: 1, Character: 33},
			{Type: RCurly, Val: "}", Line: 1, Character: 35},
			{Type: EOF},
		},
	},
	{
		input: "return .38\ntime write",
		tokens: []Token{
			{Type: Return, Val: "return", Line: 1},
			{Type: FloatLiteral, Val: ".38", Line: 1, Character: 7},
			{Type: NewLine, Val: "\n", Line: 1, Character: 10},
			{Type: Time, Val: "time", Line: 2},
			{Type: Write, Val: "write", Line: 2, Character: 5},
			{Type: EOF},
		},
	},
}

func TestLexer(t *testing.T) {
	for _, test := range tests {
		tokens, _ := Lex(test.input)
		for _, tok := range tokens {
			fmt.Println(tok.DumpString())
		}
		assert.Equal(t, test.tokens, tokens)
		fmt.Println("-------------")
	}
}


func TestLexerWithAssignments(t *testing.T) {
	tests := map[string]struct{ input, expected string }{}
	err := filepath.Walk("../../assignment1/lexer-tests1/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileName := info.Name()
		switch {
		case strings.Contains(fileName, "my-output"), info.IsDir():
			return nil
		case strings.Contains(fileName, "output"):
			key := fileName[:len(".output")]
			b, err := ioutil.ReadFile(path)
			assert.Nil(t, err)
			tests[key] = struct{ input, expected string }{input: tests[key].input, expected: string(b)}
		default:
			b, err := ioutil.ReadFile(path)
			assert.Nil(t, err)
			tests[fileName] = struct{ input, expected string }{input: string(b), expected: tests[fileName].input}
		}

		return nil
	})

	for _, test := range tests {
		buf := bytes.NewBufferString("")
		tokens, ok := Lex(test.input)
		assert.True(t, ok)
		for _, tok := range tokens {
			buf.WriteString(tok.DumpString() + "\n")
		}
		buf.WriteString("Compilation succeeded\n")
		assert.Equal(t, test.expected, buf.String())
	}
	assert.Nil(t, err)
}
