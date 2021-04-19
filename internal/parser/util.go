package parser

import "github.com/ryan-berger/jpl/internal/lexer"

// parseList takes an end character and calls parseFn to parse items between commas
func (p *parser) parseList(end lexer.TokenType, parseFn func() error) error {
	if p.expectPeek(end) {
		return nil
	}

	p.advance()

	for i := 0; i < 512; i++ {
		if err := parseFn(); err != nil {
			return err
		}

		if p.peekTokenIs(end) {
			break
		}

		if !p.expectPeek(lexer.Comma) {
			return p.errorf(p.peek, "expected ',' received %s", p.peek.Val)
		}
		p.advance()
	}
	if !p.expectPeek(end) {
		return p.errorf(p.peek, "found unexpected token or >64 elements at expression")
	}

	return nil
}
