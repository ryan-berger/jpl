package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)

func parse(t *testing.T, expr string) ast.Expression {
	tokens, ok := lexer.
		NewLexer(fmt.Sprintf("return %s", expr)).
		LexAll()

	assert.True(t, ok, "lexer error")

	commands, err := parser.
		NewParser(tokens).
		ParseProgram(false)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.ReturnStatement{}, commands[0])

	return commands[0].(*ast.ReturnStatement).Expr
}

var okTests = []struct {
	expr string
	typ  Type
}{
	{"1 + 1", Integer},
	{"1.02 + 3.2", Float},
	{"1.02 > 3.2", Boolean},
	{"1 < 3", Boolean},
	{"!true", Boolean},
	{"-1", Integer},
	{"-.023", Float},
	{"(1 + 2) / 3", Integer},
	{"(1 + 2) / 3", Integer},
	{"sum[i : 10] i + 3", Integer},
	{"if true then 10 else 22", Integer},
	{"if true then 10.11 else 22.23", Float},
	{"add_ints(10, 11)", Integer},
	{"add_Floats(1.0, 1.1)", Float},
	{"{1, 2.0}", &Tuple{Types: []Type{Integer, Float}}},
	{"{1, 2.0}{1}", Float},
	{"[[1, 2, 3, 4]][0][0]", Integer},
}

func TestExpressionCheck(t *testing.T) {
	for _, test := range okTests {
		expr := parse(t, test.expr)
		typ, err := ExpressionType(expr, NewSymbolTable())
		assert.Nil(t, err)
		assert.True(t, test.typ.Equal(typ))
	}
}

var equalityTest = []struct {
	is, other Type
	ok        bool
}{
	{Integer, Integer, true},
	{Float, Float, true},
	{Boolean, Boolean, true},
	{Integer, Float, false},
	{Boolean, Float, false},
	{Integer, Boolean, false},
	{&Array{Inner: Float, Rank: 1}, &Array{Inner: Float, Rank: 1}, true},
	{&Array{Inner: Float, Rank: 1}, &Array{Inner: Integer, Rank: 1}, false},
	{
		is:    &Array{&Array{Integer, 2}, 1},
		other: &Array{&Array{Integer, 2}, 1},
		ok:    true,
	},
	{
		is:    &Tuple{[]Type{Integer, Float, Integer}},
		other: &Tuple{[]Type{Integer, Float, Integer}},
		ok:    true,
	},
	{
		is:    &Tuple{[]Type{Integer, Float, Integer}},
		other: &Tuple{[]Type{Integer, Integer}},
		ok:    false,
	},
	{
		is:    &Tuple{[]Type{Float, Integer}},
		other: &Tuple{[]Type{Integer, Integer}},
		ok:    false,
	},
	{
		is:    &Tuple{[]Type{Integer, Integer, &Tuple{[]Type{Integer}}}},
		other: &Tuple{[]Type{Integer, Integer, &Tuple{[]Type{Integer}}}},
		ok:    true,
	},
}

func TestTypeEqual(t *testing.T) {
	for _, test := range equalityTest {
		assert.Equal(t, test.is.Equal(test.other), test.ok)
	}
}
