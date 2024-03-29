package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/lexer"
)

var tests = []struct {
	expr     string
	expected string
}{
	{
		"-a * b",
		"((-a) * b)",
	},
	{
		"!-a",
		"(!(-a))",
	},
	{
		"a + b + c",
		"((a + b) + c)",
	},
	{
		"a + b - c",
		"((a + b) - c)",
	},
	{
		"a * b * c",
		"((a * b) * c)",
	},
	{
		"a * b / c",
		"((a * b) / c)",
	},
	{
		"a + b / c",
		"(a + (b / c))",
	},
	{
		"a + b * c + d / e - f",
		"(((a + (b * c)) + (d / e)) - f)",
	},
	{
		"5 > 4 == 3 < 4",
		"((5 > 4) == (3 < 4))",
	},
	{
		"5 < 4 != 3 > 4",
		"((5 < 4) != (3 > 4))",
	},
	{
		"3 + 4 * 5 == 3 * 1 + 4 * 5",
		"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
	},
	{
		"true",
		"true",
	},
	{
		"false",
		"false",
	},
	{
		"3 > 5 == false",
		"((3 > 5) == false)",
	},
	{
		"3 < 5 == true",
		"((3 < 5) == true)",
	},
	{
		"1 + (2 + 3) + 4",
		"((1 + (2 + 3)) + 4)",
	},
	{
		"(5 + 5) * 2",
		"((5 + 5) * 2)",
	},
	{
		"2 / (5 + 5)",
		"(2 / (5 + 5))",
	},
	{
		"(5 + 5) * 2 * (5 + 5)",
		"(((5 + 5) * 2) * (5 + 5))",
	},
	{
		"-(5 + 5)",
		"(-(5 + 5))",
	},
	{
		"!(true == true)",
		"(!(true == true))",
	},
	{
		"true || false && false",
		"if true then true else if false then false else false",
	},
	{
		"0 <= x && x <= 1",
		"if (0 <= x) then (x <= 1) else false",
	},
	{
		"if x == 0 then 20 + 10 else 30",
		"if (x == 0) then (20 + 10) else 30",
	},
	{
		"sum[i : 10, j: 30 + 10 / 2] i",
		"sum[i : 10, j : (30 + (10 / 2))] i",
	},
	{
		"array[i : 4 + 3 % 30, j: 10 / 10 + 3 * 10] i",
		"array[i : (4 + (3 % 30)), j : ((10 / 10) + (3 * 10))] i",
	},
	{
		"1.234 + 2.111",
		"(1.234000 + 2.111000)",
	},
	{
		"3 + arr[10 / 2] * 2",
		"(3 + (arr[(10 / 2)] * 2))",
	},
	{
		"3 + call(3, 4, 5, 6) / 2",
		"(3 + (call(3, 4, 5, 6) / 2))",
	},
	{
		"[1, 2, 3, 4, 5][0] + 3 / 2",
		"([1, 2, 3, 4, 5][0] + (3 / 2))",
	},
	{
		"[[1, 3, 4, 5], 2, 3, 4, 5][0][1] + 3 / 2",
		"([[1, 3, 4, 5], 2, 3, 4, 5][0][1] + (3 / 2))",
	},
	{
		"{{1, 2}, 3, 4, 5}{0}{1} + 3 / 2",
		"({{1, 2}, 3, 4, 5}{0}{1} + (3 / 2))",
	},
	{
		"array[i : 5] i + array[i : 4] i",
		"array[i : 5] (i + array[i : 4] i)",
	},
	{
		"array[] sum[] array[] sum[] 1",
		"array[] sum[] array[] sum[] 1",
	},
}

func TestPrecedenceParsing(t *testing.T) {
	for _, test := range tests {
		tokens, ok := lexer.Lex(test.expr)
		assert.True(t, ok)
		parser := newParser(tokens)
		expr, err := parser.parseExpression(lowest) // parse expression
		assert.Nil(t, err)
		assert.NotNil(t, expr)
		assert.Equal(t, test.expected, expr.String())
	}
}

var errorTests = []struct {
	expr     string
	expected string
}{
	{"", "unable to parse prefix operator"},
	{"if )", "unable to parse prefix operator )"},
	{"if 3 + 3 th", "expected 'then' received 'th'"},
	{"if 4 + 3 then 30 els 10", "expected 'else' received 'els'"},
	{"[1, 2, 3, 4][]", "expected expression, found ']'"},
	{"11111111111111111111111111111111111111111111111", "integer literal 11111111111111111111111111111111111111111111111 too large for a 64 bit integer"},
	{"(3 + 4\n", "illegal token. Expected ')', found \n"},
}

func TestParseErrors(t *testing.T) {
	for _, test := range errorTests {
		tokens, ok := lexer.Lex(test.expr)
		assert.True(t, ok)
		parser := newParser(tokens)
		_, err := parser.parseExpression(lowest) // parse expression
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), test.expected)
	}

}
