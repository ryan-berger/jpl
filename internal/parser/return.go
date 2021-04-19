package parser

import "github.com/ryan-berger/jpl/internal/ast"

func  (p *parser) parseReturnStatement() (ast.Statement, error) {
	ret := &ast.ReturnStatement{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()
	var err error
	if ret.Expr, err = p.parseExpression(lowest); ret.Expr == nil {
		return nil, err
	}

	p.advance()
	return ret, nil
}
