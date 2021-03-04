package typed

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/ryan-berger/jpl/internal/symbol"
	"github.com/ryan-berger/jpl/internal/types"
)

func parseLet(t *testing.T, expr string) *ast.LetStatement {
	tokens, ok := lexer.Lex(fmt.Sprintf("let %s", expr))

	assert.True(t, ok, "lexer error")

	commands, err := parser.Parse(tokens)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.LetStatement{}, commands[0])

	return commands[0].(*ast.LetStatement)
}

var validLetTests = []struct {
	expr string
}{
	{
		expr: "x = 10",
	},
	{
		expr: "x = {1, 2}",
	},
	{
		expr: "{x, y} = {1, 2}",
	},
	{
		expr: "{x, y} = {[1, 2], 2}",
	},
	{
		expr: "{x[L], y} = {[1, 2], 2}",
	},
}

func TestLet(t *testing.T) {
	for _, test := range validLetTests {
		table := symbol.NewSymbolTable()
		letStmt := parseLet(t, test.expr)
		rType, err := expressionType(letStmt.Expr, table)
		assert.Nil(t, err)
		err = bindLVal(letStmt.LValue, rType, table)
		assert.Nil(t, err)
	}
}

func parseFunction(t *testing.T, expr string) *ast.Function {
	tokens, ok := lexer.Lex(expr)

	assert.True(t, ok, "lexer error")

	commands, err := parser.Parse(tokens)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.Function{}, commands[0])

	return commands[0].(*ast.Function)
}

var validFunctionTests = []struct {
	expr       string
	bindings   []types.Type
	returnType types.Type
}{
	{
		expr: `fn test() : {} {
}`,
		bindings:   []types.Type{},
		returnType: &types.Tuple{},
	},
	{
		expr: `fn test() : int {
return 0
}`,
		bindings:   []types.Type{},
		returnType: types.Integer,
	},
	{
		expr: `fn test() : {int, int} {
return {1, 2}
}`,
		bindings:   []types.Type{},
		returnType: &types.Tuple{Types: []types.Type{types.Integer, types.Integer}},
	},
	{
		expr: `fn test({x : int, y[H, W] : int[,]}, z : int) : {int, int} {
return {1, 2}
}`,
		bindings: []types.Type{
			&types.Tuple{
				Types: []types.Type{
					types.Integer,
					&types.Array{Inner: types.Integer, Rank: 2},
				},
			},
			types.Integer,
		},
		returnType: &types.Tuple{Types: []types.Type{types.Integer, types.Integer}},
	},
	{
		expr: `fn test({x : int, y[H, W] : int[,]}) : {int, int[]} {
return {1, [1, 2]}
}`,
		bindings: []types.Type{
			&types.Tuple{
				Types: []types.Type{
					types.Integer,
					&types.Array{Inner: types.Integer, Rank: 2},
				},
			},
		},
		returnType: &types.Tuple{Types: []types.Type{types.Integer, &types.Array{Inner: types.Integer, Rank: 1}}},
	},
	{
		expr: `fn test(x : {float, bool}) : {int, int} {
return {1, 1}
}`,
		bindings: []types.Type{
			&types.Tuple{
				Types: []types.Type{
					types.Float,
					types.Boolean,
				},
			},
		},
		returnType: &types.Tuple{Types: []types.Type{types.Integer, types.Integer}},
	},
}

func TestFunction(t *testing.T) {
	for _, test := range validFunctionTests {
		table := symbol.NewSymbolTable()
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
