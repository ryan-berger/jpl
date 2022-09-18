package parser

import (
	"testing"

	"github.com/ryan-berger/jpl/internal/lexer"
)

func FuzzParse(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		ts, _ := lexer.Lex(s)
		Parse(ts)
	})
}
