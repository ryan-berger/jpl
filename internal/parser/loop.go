package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *parser) parseArrayTransform() (ast.Expression, error) {
	expr := &ast.ArrayTransform{}
	if !p.expectPeek(lexer.LBrace) {
		return nil, p.errorf(p.peek, "expected '[' received %s", p.peek.Val)
	}

	var err error
	expr.OpBindings, err = p.parseOpBindings()
	if err !=  nil {
		return nil, err
	}

	p.advance() // move onto start of expression
	if expr.Expr, err = p.parseExpression(lowest); err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *parser) parseSumTransform() (ast.Expression, error) {
	expr := &ast.SumTransform{}
	if !p.expectPeek(lexer.LBrace) {
		return nil, p.errorf(p.peek, "expected '[' received %s", p.peek.Val)
	}

	var err error
	if expr.OpBindings, err = p.parseOpBindings(); err != nil {
		return nil, err
	}

	p.advance() // move onto start of expression

	if expr.Expr, err = p.parseExpression(array); expr.Expr == nil {
		return nil, err
	}

	return expr, nil
}

func (p *parser) parseOpBindings() ([]ast.OpBinding, error) {
	var opBindings []ast.OpBinding

	listErr := p.parseList(lexer.RBrace, func() error {
		var opBinding ast.OpBinding
		if !p.curTokenIs(lexer.Variable) {
			return p.errorf(p.peek,"expecting variable, received %s", p.peek.Val)
		}

		opBinding.Variable = p.cur.Val

		if !p.expectPeek(lexer.Colon) {
			return p.errorf(p.peek, "expecting ':', received %s", p.peek.Val)
		}
		p.advance()

		var err error
		opBinding.Expr, err = p.parseExpression(lowest)
		if err != nil {
			return err
		}

		opBindings = append(opBindings, opBinding)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}
	return opBindings, nil
}
