package parser

import "github.com/ryan-berger/jpl/internal/ast"

func  (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	ret := &ast.ReturnStatement{}
	if ret.Expr = p.parseExpression(); ret.Expr == nil {
		return nil
	}
	return ret
}
