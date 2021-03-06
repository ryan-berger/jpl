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
	or
	and
	equal
	lg
	sum
	product
	prefix
	call
)

var opPrecedence = map[lexer.TokenType]precedence{
	lexer.Or:                 or,
	lexer.And:                and,
	lexer.EqualTo:            equal,
	lexer.NotEqualTo:         equal,
	lexer.LessThan:           lg,
	lexer.LessThanOrEqual:    lg,
	lexer.GreaterThan:        lg,
	lexer.GreaterThanOrEqual: lg,
	lexer.Plus:               sum,
	lexer.Minus:              sum,
	lexer.Multiply:           product,
	lexer.Divide:             product,
	lexer.Mod:                product,
	lexer.LCurly:             call,
	lexer.LBrace:             call,
}

func (p *parser) parseExpression(pr precedence) ast.Expression {
	prefix := p.prefixParseFns[p.cur.Type]
	if prefix == nil {
		p.errorf(p.cur, "unable to parse prefix operator %s", p.cur.Val)
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

func (p *parser) parsePrefixExpr() ast.Expression {
	expr := &ast.PrefixExpression{
		Op: p.cur.Val,
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()
	expr.Expr = p.parseExpression(prefix)
	if expr.Expr == nil {
		return nil
	}

	return expr
}

func (p *parser) parseInfixExpr(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Op:   p.cur.Val,
		Left: left,
	}
	pr := p.curPrecedence()
	p.advance()

	expr.Right = p.parseExpression(pr)
	return expr
}

func (p *parser) parseArrayRefExpr(arr ast.Expression) ast.Expression {
	arrRefExpr := &ast.ArrayRefExpression{
		Array: arr,
	}
	ok := p.parseList(lexer.RBrace, func() bool {
		expr := p.parseExpression(lowest)
		if expr == nil {
			return false
		}
		arrRefExpr.Indexes = append(arrRefExpr.Indexes, expr)
		return true
	})

	if !ok {
		return nil
	}

	if len(arrRefExpr.Indexes) == 0 {
		p.errorf(p.cur, "expected expression, found ']'")
		return nil
	}

	return arrRefExpr
}

func (p *parser) parseTupleRefExpr(tuple ast.Expression) ast.Expression {
	arrRefExpr := &ast.TupleRefExpression{
		Tuple: tuple,
	}
	p.advance()
	if arrRefExpr.Index = p.parseExpression(lowest); arrRefExpr.Index == nil {
		return nil
	}
	p.advance()
	return arrRefExpr
}

func (p *parser) parseGroupedExpression() ast.Expression {
	p.advance()

	exp := p.parseExpression(lowest) // TODO: handle error

	if exp == nil {
		return nil
	}

	if !p.expectPeek(lexer.RParen) {
		p.errorf(p.peek, "illegal token. Expected ')', found %s", p.peek.Val)
		return nil
	}

	return exp
}

func (p *parser) parseTupleExpression() ast.Expression {
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

func (p *parser) parseArrayExpression() ast.Expression {
	arrayExpr := &ast.ArrayExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	ok := p.parseList(lexer.RBrace, func() bool {
		expr := p.parseExpression(lowest)
		if expr == nil {
			return false
		}

		arrayExpr.Expressions = append(arrayExpr.Expressions, expr)
		return true
	})

	if !ok {
		return nil
	}

	return arrayExpr
}

func (p *parser) parseInteger() ast.Expression {
	expr := &ast.IntExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	val, err := strconv.ParseInt(p.cur.Val, 10, 64)
	if err != nil {
		p.errorf(p.cur, "integer literal %s too large for a 64 bit integer", p.cur.Val)
		return nil
	}

	expr.Val = val
	return expr
}

func (p *parser) parseFloat() ast.Expression {
	expr := &ast.FloatExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	val, err := strconv.ParseFloat(p.cur.Val, 64)
	if err != nil {
		p.errorf(p.cur, "float %s too large for a 64 bit float", p.cur.Val)
		return nil
	}
	expr.Val = val
	return expr
}

func (p *parser) parseBoolean() ast.Expression {
	return &ast.BooleanExpression{
		Val: p.cur.Val == "true",
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
}

func (p *parser) parseIdentifier() ast.Expression {
	val := p.cur.Val
	if p.peekTokenIs(lexer.LParen) {
		return p.parseCallExpression()
	}
	return &ast.IdentifierExpression{
		Identifier: val,
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
}

func (p *parser) parseCallExpression() ast.Expression {
	callExpr := &ast.CallExpression{
		Identifier: p.cur.Val,
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	ok := p.parseList(lexer.RParen, func() bool {
		expr := p.parseExpression(lowest)
		if expr == nil {
			return false
		}
		callExpr.Arguments = append(callExpr.Arguments, expr)
		return true
	})

	if !ok {
		return nil
	}
	return callExpr
}

func (p *parser) parseIf() ast.Expression {
	expr := &ast.IfExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	if expr.Condition = p.parseExpression(lowest); expr.Condition == nil {
		return nil
	}

	if p.error != nil {
		return nil
	}

	if !p.expectPeek(lexer.Then) {
		p.errorf(p.peek, "expected 'then' received '%s'", p.peek.Val)
		return nil
	}
	p.advance()

	if expr.Consequence = p.parseExpression(lowest); expr.Consequence == nil {
		return nil
	}

	if p.error != nil {
		return nil
	}

	if !p.expectPeek(lexer.Else) {
		p.errorf(p.peek, "expected 'else' received '%s'", p.peek.Val)
		return nil
	}
	p.advance()

	if expr.Otherwise = p.parseExpression(lowest); expr.Otherwise == nil {
		return nil
	}

	return expr
}
