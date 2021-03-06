package parser

import "github.com/ryan-berger/jpl/internal/lexer"

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
			p.errorf("error expected ',' received %s at %d:%d", p.peek.Val, p.peek.Line, p.peek.Character)
			return false
		}
		p.advance()
	}
	if !p.expectPeek(end) {
		p.errorf("error found unexpected token or >64 elements at expression at %d:%d", p.peek.Line, p.peek.Character)
		return false
	}

	return true
}
