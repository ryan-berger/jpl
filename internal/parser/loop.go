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

	p.advance() // move onto start of expression

	if expr.Expr = p.parseExpression(lowest); expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *Parser) parseOpBindings() []ast.OpBinding {
	var opBindings []ast.OpBinding

	ok := p.parseList(lexer.RBrace, func() bool {
		var opBinding ast.OpBinding
		if !p.curTokenIs(lexer.Variable) {
			p.errorf("expecting variable, received %s at line %d", p.peek.Val, p.peek.Line)
			return false
		}

		opBinding.Variable = p.cur.Val

		if !p.expectPeek(lexer.Colon) {
			p.errorf("expecting ':', received %s at line %d", p.peek.Val, p.peek.Line)
			return false
		}
		p.advance()

		opBinding.Expr = p.parseExpression(lowest)
		if opBinding.Expr == nil {
			return false
		}

		opBindings = append(opBindings, opBinding)
		return true
	})

	if !ok {
		return nil
	}
	return opBindings
}
