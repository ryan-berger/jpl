package checker

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

func parse(t *testing.T, expr string) ast.Expression {
	tokens, ok := lexer.
		Lex(fmt.Sprintf("return %s", expr))

	assert.True(t, ok, "lexer error")

	commands, err := parser.Parse(tokens)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.ReturnStatement{}, commands[0])

	return commands[0].(*ast.ReturnStatement).Expr
}

var okTests = []struct {
	expr string
	typ  types.Type
}{
	{"1 + 1", types.Integer},
	{"1.02 + 3.2", types.Float},
	{"1.02 > 3.2", types.Boolean},
	{"1 < 3", types.Boolean},
	{"!true", types.Boolean},
	{"-1", types.Integer},
	{"-.023", types.Float},
	{"(1 + 2) / 3", types.Integer},
	{"(1 + 2) / 3", types.Integer},
	{"sum[i : 10] i + 3", types.Integer},
	{"array[i : 10] [i]", &types.Array{Inner: types.Integer, Rank: 1}},
	{"if true then 10 else 22", types.Integer},
	{"if true then 10.11 else 22.23", types.Float},
	{"sub_ints(10, 11)", types.Integer},
	{"sub_floats(1.0, 1.1)", types.Float},
	{"{1, 2.0}", &types.Tuple{Types: []types.Type{types.Integer, types.Float}}},
	{"{1, 2.0}{1}", types.Float},
	{"[[1, 2, 3, 4]][0][0]", types.Integer},
	{"[1, 2, 3, 4]", &types.Array{Inner: types.Integer, Rank: 1}},
	{"[]", &types.Array{Inner: types.Integer, Rank: 1}},
}

func TestExpressionCheck(t *testing.T) {
	for _, test := range okTests {
		expr := parse(t, test.expr)
		typ, err := expressionType(expr, symbol.NewSymbolTable())
		assert.Nil(t, err)
		assert.True(t, test.typ.Equal(typ))
	}
}

var equalityTest = []struct {
	is, other types.Type
	ok        bool
}{
	{types.Integer, types.Integer, true},
	{types.Float, types.Float, true},
	{types.Boolean, types.Boolean, true},
	{types.Integer, types.Float, false},
	{types.Boolean, types.Float, false},
	{types.Integer, types.Boolean, false},
	{&types.Array{Inner: types.Float, Rank: 1}, &types.Array{Inner: types.Float, Rank: 1}, true},
	{&types.Array{Inner: types.Float, Rank: 1}, &types.Array{Inner: types.Integer, Rank: 1}, false},
	{
		is:    &types.Array{&types.Array{types.Integer, 2}, 1},
		other: &types.Array{&types.Array{types.Integer, 2}, 1},
		ok:    true,
	},
	{
		is:    &types.Tuple{[]types.Type{types.Integer, types.Float, types.Integer}},
		other: &types.Tuple{[]types.Type{types.Integer, types.Float, types.Integer}},
		ok:    true,
	},
	{
		is:    &types.Tuple{[]types.Type{types.Integer, types.Float, types.Integer}},
		other: &types.Tuple{[]types.Type{types.Integer, types.Integer}},
		ok:    false,
	},
	{
		is:    &types.Tuple{[]types.Type{types.Float, types.Integer}},
		other: &types.Tuple{[]types.Type{types.Integer, types.Integer}},
		ok:    false,
	},
	{
		is:    &types.Tuple{[]types.Type{types.Integer, types.Integer, &types.Tuple{[]types.Type{types.Integer}}}},
		other: &types.Tuple{[]types.Type{types.Integer, types.Integer, &types.Tuple{[]types.Type{types.Integer}}}},
		ok:    true,
	},
}

func TestTypeEqual(t *testing.T) {
	for _, test := range equalityTest {
		assert.Equal(t, test.is.Equal(test.other), test.ok)
	}
}

var failureTests = []struct{
	expr string
	err string
}{
	{"if true then 0 else 2.0", "branches return different types: int, float"},
	{"if x && y then 0 else 2.0", "unknown symbol x"},
	{"if true then x + y else 2.0", "unknown symbol x"},
	{"if true then 2 else x + y", "unknown symbol x"},
	{"if 22 then 0 else 3", "expected boolean, received int at position 1:10"},
	{"sub_floats + 2", "found function name, expected identifier sub_floats"},
	{"sub_float(1, 2)", "unknown symbol sub_float"},
	{"sub_floats(1.0, 2.0, 3.4)", "function sub_floats expects 2 arguments, received 3"},
	{"sub_floats(1.0, 2)", "type error: expected float received int"},
	{"sub_floats(1, 2.0)", "type error: expected float received int"},
	{"sub_floats(x, 2)", "unknown symbol x"},
	{"1 && true", "type error: left operand of && expression is of type int expected bool"},
	{"true && 1", "type error: right operand of && expression is of type int expected bool"},
	{"true + 1", "type error: left type of + expression must be numeric, received bool"},
	{"1 + true", "type error: right type of + expression must be numeric, received bool"},
	{"1 + 1.0", "type error: both sides of numerical operation must be of the same type"},
	{"1 + x", "unknown symbol x"},
	{"!1", "type error, expected boolean on right hand side of '!'"},
	{"!x", "unknown symbol x"},
	{"-[1, 2, 3]", "type error, expected numeric type on right hand side of '-'"},
	{"[1, 2, 3]{1}", "tuple index of non-tuple type"},
	{"{1, 2, 3}{1+2}", "tuple indexing requires integer literal"},
	{"{1, 2, 3}{4}", "tuple index out of bounds"},
	{"{1, x, 3}{2}", "unknown symbol x"},
	{"[1, x, 3][0]", "unknown symbol x"},
	{"{1, 2, 3}[0]", "array reference of non-array"},
	{"[1, 2, 3][0, 1]", "array access of rank 1 with 2 indexes"},
	{"[1, 2, 3][x]", "unknown symbol x"},
	{"[1, 2, 3][2.0]", "non-integer index of array type float"},
	{"[[1, 2], 2, 3]", "array literal has mixed types"},
	{"[x, 1, 3]", "unknown symbol x"},
	{"sum[i : 10] j + 10", "unknown symbol j"},
	{"sum[i : 10, i : 10] j + 10", "illegal shadowing in sum expr, var: i"},
	{"sum[i : 10] true && false", "sum returns non-numeric expression"},
	{"sum[i : x + 10] 1", "unknown symbol x"},
	{"sum[i : true] 1", "bindArg expr initializer for i returns non-integer"},
	{"array[i : 10] j + 10", "unknown symbol j"},
	{"array[i : 10, i : 10] j + 10", "illegal shadowing in sum expr, var: i"},
	{"array[i : 10, j: 10] [i + j]", "return type of array expression must be of equal rank of number of bindings"},
	{"array[i : 10] i", "return type of array expression must be array"},
	{"array[i : x + 10] 1", "unknown symbol x"},
	{"array[i : true] 1", "bindArg expr initializer for i returns non-integer"},
}

func TestCheckFailures(t *testing.T) {
	table := symbol.NewSymbolTable()
	for _, test := range failureTests {
		expr := parse(t, test.expr)
		_, err := expressionType(expr, table)
		assert.NotNil(t, err, test.expr)
		assert.Contains(t, err.Error(), test.err)
	}
}

