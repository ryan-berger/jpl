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
		p.errorf("err: illegal token. Expected string, found %s at line %d", p.peek.Val, p.peek.Line)
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

	if p.expectPeek(lexer.RCurly) {
		p.advance() // move past curly
		return tupleExpr
	}

	p.advance() // move past lBrace

	tupleExpr.Expressions = append(tupleExpr.Expressions, p.parseExpression(lowest)) // TODO: error handling
	for p.peekTokenIs(lexer.Comma) {
		p.advance()
		p.advance()
		tupleExpr.Expressions = append(tupleExpr.Expressions, p.parseExpression(lowest))
	}

	if !p.expectPeek(lexer.RCurly) {
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
	return &ast.IdentifierExpression{Identifier: p.cur.Val}
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
