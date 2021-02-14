package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/lexer"
)

func TestParser_advance(t *testing.T) {
	tokens := []lexer.Token{
		{Type: lexer.And, Val: "&&"},
		{Type: lexer.Or, Val: "||"},
		{Type: lexer.Not, Val: "!"},
	}
	p := NewParser(tokens)
	assert.Equal(t, p.tokens[0], p.cur)
	assert.Equal(t, p.tokens[1], p.peek)
}
