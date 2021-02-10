package parser

import "github.com/ryan-berger/jpl/internal/lexer"

func (p *Parser) parseList(end lexer.TokenType, parseFn func() bool) bool {
	if p.expectPeek(end) {
		p.advance()
		return true
	}

	p.advance()

	for i := 0; i < 64; i++ {
		if !parseFn() {
			return false
		}

		if p.peekTokenIs(end) {
			break
		}

		if !p.expectPeek(lexer.Comma) {
			p.errorf("error expected ',' received %s at line %d", p.peek.Val, p.peek.Line)
			return false
		}
		p.advance()
	}
	if !p.expectPeek(end) {
		p.errorf("error found unexpected token or >64 elements at expression at line %d", p.peek.Line)
	}

	return true
}
