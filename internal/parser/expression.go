package parser

import (
	"strconv"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type precedence int

const (
	_ precedence = iota
	lowest
	equals
	lg
	sum
	product
	prefix
	call
)

var opPrecedence = map[lexer.TokenType]precedence{
	lexer.EqualTo:     equals,
	lexer.NotEqualTo:  equals,
	lexer.LessThan:    lg,
	lexer.GreaterThan: lg,
	lexer.Plus:        sum,
	lexer.Minus:       sum,
	lexer.Or:          sum,
	lexer.Multiply:    product,
	lexer.Divide:      product,
	lexer.Mod:         product,
	lexer.And:         product,
}

func (p *Parser) parseExpression(pr precedence) ast.Expression {
	prefix := p.prefixParseFns[p.cur.Type]
	if prefix == nil { // TODO: actually implement
		return nil
	}

	leftExp := prefix()
	for (!p.peekTokenIs(lexer.NewLine) || !p.peekTokenIs(lexer.EOF)) && pr < p.peekPrecedence() {
		infix := p.infixParseFns[p.peek.Type]
		if infix == nil {
			return leftExp
		}
		p.advance()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parsePrefixExpr() ast.Expression {
	expr := &ast.PrefixExpression{
		Op: p.cur.Val,
	}
	p.advance()
	expr.Expr = p.parseExpression(prefix)
	if expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Op:   p.cur.Val,
		Left: left,
	}
	pr := p.curPrecedence()
	p.advance()

	expr.Right = p.parseExpression(pr)
	return expr
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advance()

	exp := p.parseExpression(lowest) // TODO: handle error

	if exp == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		p.errorf("err: illegal token. Expected ), found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	return exp
}

func (p *Parser) parseTupleExpression() ast.Expression {
	tupleExpr := &ast.TupleExpression{}

	ok := p.parseList(lexer.RCurly, func() bool {
		expr := p.parseExpression(lowest)
		if expr == nil {
			return false
		}

		tupleExpr.Expressions = append(tupleExpr.Expressions, expr)
		return true
	})

	if !ok {
		return nil
	}
	return tupleExpr
}

func (p *Parser) parseInteger() ast.Expression {
	expr := &ast.IntExpression{}
	val, err := strconv.ParseInt(p.cur.Val, 10, 64)
	if err != nil {
		return nil
	}

	expr.Val = val
	return expr
}

func (p *Parser) parseFloat() ast.Expression {
	expr := &ast.FloatExpression{}
	val, err := strconv.ParseFloat(p.cur.Val, 64)
	if err != nil {
		return nil
	}
	expr.Val = val
	return expr
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanExpression{
		Val: p.cur.Val == "true",
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	val := p.cur.Val
	if !p.expectPeek(lexer.LParen) {
		return &ast.IdentifierExpression{Identifier: val}
	}

	var exprs []ast.Expression
	ok := p.parseList(lexer.RParen, func() bool {
		expr := p.parseExpression(lowest)
		if expr == nil {
			return false
		}
		exprs = append(exprs, expr)
		return true
	})

	if !ok {
		return nil
	}
	return &ast.CallExpression{Identifier: val, Arguments: exprs}
}

func (p *Parser) parseIf() ast.Expression {
	expr := &ast.IfExpression{}
	p.advance()

	if expr.Condition = p.parseExpression(lowest); expr.Condition == nil {
		return nil
	}

	p.expectPeek(lexer.Then)
	p.advance()

	if expr.Consequence = p.parseExpression(lowest); expr.Consequence == nil {
		return nil
	}

	p.expectPeek(lexer.Else)
	p.advance()

	if expr.Otherwise = p.parseExpression(lowest); expr.Otherwise == nil {
		return nil
	}

	return expr
}
