package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)

func parseLet(t *testing.T, expr string) ast.Expression {
	tokens, ok := lexer.
		NewLexer(fmt.Sprintf("let %s", expr)).
		LexAll()

	assert.True(t, ok, "lexer error")

	commands, err := parser.
		NewParser(tokens).
		ParseProgram(false)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.LetStatement{}, commands[0])

	return commands[0].(*ast.ReturnStatement).Expr
}

func parseFunction(t *testing.T, expr string) *ast.Function {
	tokens, ok := lexer.
		NewLexer(expr).
		LexAll()

	assert.True(t, ok, "lexer error")

	commands, err := parser.
		NewParser(tokens).
		ParseProgram(false)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.Function{}, commands[0])

	return commands[0].(*ast.Function)
}

var validFunctionTests = []struct {
	expr       string
	bindings   []Type
	returnType Type
}{
	{
		expr: `fn test() : {} {
}`,
		bindings:   []Type{},
		returnType: &Tuple{},
	},
	{
		expr: `fn test() : int {
return 0
}`,
		bindings:   []Type{},
		returnType: Integer,
	},
	{
		expr: `fn test() : {int, int} {
return {1, 2}
}`,
		bindings:   []Type{},
		returnType: &Tuple{Types: []Type{Integer, Integer}},
	},
	{
		expr: `fn test({x : int, y[H, W] : int[,]}, z : int) : {int, int} {
return {1, 2}
}`,
		bindings: []Type{
			&Tuple{
				Types: []Type{
					Integer,
					&Array{Inner: Integer, Rank: 2},
				},
			},
			Integer,
		},
		returnType: &Tuple{Types: []Type{Integer, Integer}},
	},
	{
		expr: `fn test({x : int, y[H, W] : int[,]}) : {int, int[]} {
return {1, [1, 2]}
}`,
		bindings:   []Type{
			&Tuple{
				Types: []Type{
					Integer,
					&Array{Inner: Integer, Rank: 2},
				},
			},
		},
		returnType: &Tuple{Types: []Type{Integer, &Array{Inner: Integer, Rank: 1}}},
	},
	{
		expr: `fn test(x : {float, bool}) : {int, int} {
return {1, 1}
}`,
		bindings:   []Type{
			&Tuple{
				Types: []Type{
					Float,
					Boolean,
				},
			},
		},
		returnType: &Tuple{Types: []Type{Integer, Integer}},
	},
}

func TestFunction(t *testing.T) {
	for _, test := range validFunctionTests {
		table := NewSymbolTable()
		fn := parseFunction(t, test.expr)
		symb, err := functionBinding(fn, table)

		assert.Nil(t, err)
		assert.True(t, test.returnType.Equal(symb.Return))
		assert.Len(t, symb.Args, len(test.bindings))
		for i, typ := range test.bindings {
			assert.True(t, symb.Args[i].Equal(typ))
		}
	}
}
