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
		"(true || (false && false))",
	},
	{
		expr:     "0 <= x && x <= 1",
		expected: "((0 <= x) && (x <= 1))",
	},
}

func TestPrecedenceParsing(t *testing.T) {
	for _, test := range tests {
		tokens, ok := lexer.NewLexer(test.expr).LexAll()
		assert.True(t, ok)
		parser := NewParser(tokens)
		expr := parser.parseExpression(lowest) // parse expression
		assert.Empty(t, parser.errors)
		assert.NotNil(t, expr)
		assert.Equal(t, test.expected, expr.String())
	}
}
