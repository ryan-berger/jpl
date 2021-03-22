package parser

import "github.com/ryan-berger/jpl/internal/lexer"

// parseList takes an end character and calls parseFn to parse items between commas
func (p *parser) parseList(end lexer.TokenType, parseFn func() bool) bool {
	if p.expectPeek(end) {
		return true
	}

	p.advance()

	for i := 0; i < 512; i++ {
		if !parseFn() {
			return false
		}

		if p.peekTokenIs(end) {
			break
		}

		if !p.expectPeek(lexer.Comma) {
			p.errorf(p.peek, "expected ',' received %s", p.peek.Val)
			return false
		}
		p.advance()
	}
	if !p.expectPeek(end) {
		p.errorf(p.peek, "found unexpected token or >64 elements at expression")
		return false
	}

	return true
}
