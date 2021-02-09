package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseArrayTransform() ast.Expression {
	expr := &ast.ArrayTransform{}
	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	if expr.OpBindings = p.parseOpBindings(); len(expr.OpBindings) == 0 {
		return nil
	}

	if !p.expectPeek(lexer.RBrace) {
		return nil
	}

	p.advance() // move onto start of expression

	if expr.Expr = p.parseExpression(lowest); expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *Parser) parseSumTransform() ast.Expression {
	expr := &ast.SumTransform{}
	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	if expr.OpBindings = p.parseOpBindings(); len(expr.OpBindings) == 0 {
		return nil
	}

	if !p.expectPeek(lexer.RBrace) {
		return nil
	}

	p.advance() // move onto start of expression

	if expr.Expr = p.parseExpression(lowest); expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *Parser) parseOpBindings() []ast.OpBinding {
	var opBindings []ast.OpBinding

	for i := 0; i < 64; i++ {
		var opBinding ast.OpBinding
		if !p.expectPeek(lexer.Variable) {
			return nil
		}
		opBinding.Variable = p.cur.Val

		if !p.expectPeek(lexer.Colon) {
			return nil
		}
		p.advance()

		opBinding.Expr = p.parseExpression(lowest)
		if opBinding.Expr == nil {
			return nil
		}

		opBindings = append(opBindings, opBinding)
		if !p.expectPeek(lexer.Comma) {
			break
		}
	}
	return opBindings
}