package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *parser) parseArrayTransform() ast.Expression {
	expr := &ast.ArrayTransform{}
	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	expr.OpBindings = p.parseOpBindings()
	if p.error !=  nil {
		return nil
	}

	p.advance() // move onto start of expression
	if expr.Expr = p.parseExpression(lowest); expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *parser) parseSumTransform() ast.Expression {
	expr := &ast.SumTransform{}
	if !p.expectPeek(lexer.LBrace) {
		return nil
	}

	if expr.OpBindings = p.parseOpBindings(); len(expr.OpBindings) == 0 {
		return nil
	}

	p.advance() // move onto start of expression

	if expr.Expr = p.parseExpression(array); expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *parser) parseOpBindings() []ast.OpBinding {
	var opBindings []ast.OpBinding

	ok := p.parseList(lexer.RBrace, func() bool {
		var opBinding ast.OpBinding
		if !p.curTokenIs(lexer.Variable) {
			p.errorf(p.peek,"expecting variable, received %s", p.peek.Val)
			return false
		}

		opBinding.Variable = p.cur.Val

		if !p.expectPeek(lexer.Colon) {
			p.errorf(p.peek, "expecting ':', received %s", p.peek.Val)
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
