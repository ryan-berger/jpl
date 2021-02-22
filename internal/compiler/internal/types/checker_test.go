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
	{"array[i : 10] [i]", &Array{Inner: Integer, Rank: 1}},
	{"if true then 10 else 22", Integer},
	{"if true then 10.11 else 22.23", Float},
	{"add_ints(10, 11)", Integer},
	{"add_floats(1.0, 1.1)", Float},
	{"{1, 2.0}", &Tuple{Types: []Type{Integer, Float}}},
	{"{1, 2.0}{1}", Float},
	{"[[1, 2, 3, 4]][0][0]", Integer},
	{"[1, 2, 3, 4]", &Array{Inner: Integer, Rank: 1}},
	{"[]", &Array{Inner: Integer, Rank: 1}},
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

var failureTests = []struct{
	expr string
	err string
}{
	{"if true then 0 else 2.0", "branches do not return same type"},
	{"if x && y then 0 else 2.0", "unknown symbol x"},
	{"if true then x + y else 2.0", "unknown symbol x"},
	{"if true then 2 else x + y", "unknown symbol x"},
	{"if 22 then 0 else 3", "conditional expression is not boolean"},
	{"add_floats + 2", "found function name, expected identifier add_floats"},
	{"add_float(1, 2)", "unknown symbol add_float"},
	{"add_floats(1.0, 2.0, 3.4)", "function add_floats expects 2 arguments, received 3"},
	{"add_floats(1.0, 2)", "type error at arg 2"},
	{"add_floats(x, 2)", "unknown symbol x"},
	{"1 && true", "type error left operand"},
	{"true && 1", "type error right operand"},
	{"true + 1", "type error: left operand not numerical"},
	{"1 + true", "type error: right operand not numerical"},
	{"1 + 1.0", "type mismatch"},
	{"1 + x", "unknown symbol x"},
	{"!1", "type error, expected boolean on right hand side of '!'"},
	{"!x", "unknown symbol x"},
	{"-[1, 2, 3]", "type error, expected numeric type on right hand side of '-'"},
	{"[1, 2, 3]{1}", "tuple reference of non-tuple"},
	{"{1, 2, 3}{1+2}", "expected integer literal received expression"},
	{"{1, 2, 3}{4}", "tuple index out of bounds"},
	{"{1, x, 3}{2}", "unknown symbol x"},
	{"[1, x, 3][0]", "unknown symbol x"},
	{"{1, 2, 3}[0]", "array reference of non-array"},
	{"[1, 2, 3][0, 1]", "array access of rank 1 with 2 indexes"},
	{"[1, 2, 3][x]", "unknown symbol x"},
	{"[1, 2, 3][2.0]", "non-integer index expression of array"},
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
	table := NewSymbolTable()
	for _, test := range failureTests {
		expr := parse(t, test.expr)
		_, err := ExpressionType(expr, table)
		assert.NotNil(t, err, test.expr)
		assert.Equal(t, test.err, err.Error())
	}
}
