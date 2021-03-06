package expander

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
)


func parseExpression(t *testing.T, expr string) ast.Expression {
	tokens, ok := lexer.
		Lex(fmt.Sprintf("return %s", expr))

	assert.True(t, ok, "lexer error")

	commands, err := parser.Parse(tokens)

	assert.Nil(t, err, "error nil")
	assert.Len(t, commands, 1)
	assert.IsType(t, &ast.ReturnStatement{}, commands[0])

	return commands[0].(*ast.ReturnStatement).Expr
}

func newNexter() func() int {
	N := 0
	next := func() int {
		N++
		return N
	}
	return next
}

func TestExpressionExpansion(t *testing.T)  {
	exp := parseExpression(t, "resize(crop(sepia(img), 50, 250, 650, 650), 300, 200)")
	_, smts := expansionAndLet(exp, newNexter())
	for _, s := range smts {
		fmt.Println(s)
	}


}
