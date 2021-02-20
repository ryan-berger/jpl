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
		expr: `fn test() : int {
return 0
}`,
		bindings:   []Type{},
		returnType: Integer,
	},
}

func TestFunction(t *testing.T) {
	for _, test := range validFunctionTests {
		table := NewSymbolTable()
		fn := parseFunction(t, test.expr)
		symb, err := functionBinding(fn, table)

		assert.Nil(t, err)
		assert.Equal(t, test.returnType, symb.Return)
	}
}
