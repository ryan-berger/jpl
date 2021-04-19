package optimizer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/ast/dsl"
)

func TestConstantFolding(t *testing.T) {
	var tests = []struct {
		expr     ast.Expression
		expected ast.Expression
	}{
		{
			expr:     dsl.Int(5),
			expected: dsl.Int(5),
		},
		{
			expr:     dsl.Float(5),
			expected: dsl.Float(5),
		},
		{
			expr:     dsl.Bool(true),
			expected: dsl.Bool(true),
		},
		{
			expr:     dsl.Infix("+", dsl.Int(2), dsl.Int(3)),
			expected: dsl.Int(5),
		},
		{
			expr:     dsl.Infix("&&", dsl.Bool(true), dsl.Bool(false)),
			expected: dsl.Bool(false),
		},
		{
			expr:     dsl.Infix("||", dsl.Bool(true), dsl.Bool(false)),
			expected: dsl.Bool(true),
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, constantFold(test.expr))
	}
}
