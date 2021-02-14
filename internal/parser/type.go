package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseTypeExpression() ast.Type {
	t := p.parseType()

	// handle type nesting
	for p.expectPeek(lexer.LBrace) {
		rank := 1
		for p.expectPeek(lexer.Comma) {
			rank++
		}

		if !p.expectPeek(lexer.RBrace) {
			return nil
		}
		t = &ast.ArrType{
			Type: t,
			Rank: rank,
		}
	}

	return t
}

func (p *Parser) parseType() ast.Type {
	if p.curTokenIs(lexer.LCurly) {
		return p.parseTupleType()
	}
	var t ast.Type

	switch p.cur.Type {
	case lexer.Float, lexer.Int:
		t = tokenToType[p.cur.Type]
	case lexer.Float3:
		t = &ast.ArrType{
			Type: ast.Float,
			Rank: 3,
		}
	case lexer.Float4:
		t = &ast.ArrType{
			Type: ast.Float,
			Rank: 4,
		}
	default:
		p.errorf("err: expected type received %s at line %d", p.cur.Val, p.cur.Line)
		return nil
	}

	return t
}

func (p *Parser) parseTupleType() ast.Type {
	tupleType := &ast.TupleType{}

	ok := p.parseList(lexer.RCurly, func() bool {
		typ := p.parseTypeExpression()
		if typ == nil {
			return false
		}
		tupleType.Types = append(tupleType.Types, typ)
		return true
	})

	if !ok {
		return nil
	}

	return tupleType
}
