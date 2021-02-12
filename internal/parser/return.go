package parser

import "github.com/ryan-berger/jpl/internal/ast"

func  (p *Parser) parseReturnStatement() ast.Statement {
	ret := &ast.ReturnStatement{}
	p.advance()
	if ret.Expr = p.parseExpression(lowest); ret.Expr == nil {
		return nil
	}

	p.advance()
	return ret
}
