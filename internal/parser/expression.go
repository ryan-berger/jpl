package parser

import (
	"math"
	"strconv"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
)

type precedence int

const (
	_ precedence = iota
	array
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

func (p *parser) parseExpression(pr precedence) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.cur.Type]
	if prefix == nil {
		return nil, p.errorf(p.cur, "unable to parse prefix operator %s", p.cur.Val)
	}

	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}
	for (!p.peekTokenIs(lexer.NewLine) || !p.peekTokenIs(lexer.EOF)) && pr < p.peekPrecedence() {
		infix := p.infixParseFns[p.peek.Type]
		if infix == nil {
			return leftExp, nil
		}
		p.advance()
		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

const minInt = "9223372036854775808"

func (p *parser) parsePrefixExpr() (ast.Expression, error) {
	// defer untrace(trace("PREFIX"))
	expr := &ast.PrefixExpression{
		Op: p.cur.Val,
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	// MIN_INT
	if expr.Op == "-" && p.cur.Type == lexer.IntLiteral && p.cur.Val == minInt{
		return &ast.IntExpression{
			Val: math.MinInt64,
			Location: ast.Location{
				Line: p.cur.Line,
				Pos:  p.cur.Character,
			},
		}, nil
	}

	var err error
	expr.Expr, err = p.parseExpression(prefix)
	if err != nil {
		return nil, err
	}

	return expr, nil
}

func (p *parser) parseInfixExpr(left ast.Expression) (ast.Expression, error) {
	// defer untrace(trace(fmt.Sprintf("INFIX %s", p.cur.Val)))

	op := p.cur.Val

	pr := p.curPrecedence()
	p.advance()

	right, err := p.parseExpression(pr)
	if err != nil {
		return nil, err
	}

	var expr ast.Expression
	switch op {
	case "&&":
		expr = &ast.IfExpression{
			Condition:   left,
			Consequence: right,
			Otherwise:   &ast.BooleanExpression{Val: false},
		}
	case "||":
		expr = &ast.IfExpression{
			Condition:   left,
			Consequence: &ast.BooleanExpression{Val: true},
			Otherwise:   right,
		}
	default:
		expr = &ast.InfixExpression{
			Op:    op,
			Left:  left,
			Right: right,
		}
	}
	return expr, nil
}

func (p *parser) parseArrayRefExpr(arr ast.Expression) (ast.Expression, error) {

	arrRefExpr := &ast.ArrayRefExpression{
		Array: arr,
	}
	listErr := p.parseList(lexer.RBrace, func() error {
		expr, err := p.parseExpression(lowest)
		if err != nil {
			return err
		}
		arrRefExpr.Indexes = append(arrRefExpr.Indexes, expr)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}

	if len(arrRefExpr.Indexes) == 0 {
		return nil, p.errorf(p.cur, "expected expression, found ']'")
	}

	return arrRefExpr, nil
}

func (p *parser) parseTupleRefExpr(tuple ast.Expression) (ast.Expression, error) {
	tupleRefExpr := &ast.TupleRefExpression{
		Tuple: tuple,
	}

	if !p.expectPeek(lexer.IntLiteral) {
		return nil, p.errorf(p.peek,
			"illegal token, integer literal expected, found: %s", p.peek.Val)
	}

	val, err := strconv.ParseInt(p.cur.Val, 10, 64)
	if err != nil {
		return nil, p.errorf(p.cur, "integer literal %s too large for a 64 bit integer", p.cur.Val)
	}

	tupleRefExpr.Index = val

	if !p.expectPeek(lexer.RCurly) {
		return nil, p.errorf(p.peek,
			"illegal token, expected '}' received %s", p.peek.Val)
	}

	return tupleRefExpr, nil
}

func (p *parser) parseGroupedExpression() (ast.Expression, error) {
	p.advance()

	exp, err := p.parseExpression(lowest) // TODO: handle error

	if err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.RParen) {
		return nil, p.errorf(p.peek, "illegal token. Expected ')', found %s", p.peek.Val)
	}

	return exp, nil
}

func (p *parser) parseTupleExpression() (ast.Expression, error) {
	tupleExpr := &ast.TupleExpression{}

	listErr := p.parseList(lexer.RCurly, func() error {
		expr, err := p.parseExpression(lowest)
		if err != nil {
			return err
		}

		tupleExpr.Expressions = append(tupleExpr.Expressions, expr)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}
	return tupleExpr, nil
}

func (p *parser) parseArrayExpression() (ast.Expression, error) {
	arrayExpr := &ast.ArrayExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	listErr := p.parseList(lexer.RBrace, func() error {
		expr, err := p.parseExpression(lowest)
		if err != nil {
			return err
		}

		arrayExpr.Expressions = append(arrayExpr.Expressions, expr)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}

	return arrayExpr, nil
}

func (p *parser) parseInteger() (ast.Expression, error) {
	expr := &ast.IntExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	val, err := strconv.ParseInt(p.cur.Val, 10, 64)
	if err != nil {
		return nil, p.errorf(p.cur, "integer literal %s too large for a 64 bit integer", p.cur.Val)
	}

	expr.Val = val
	return expr, nil
}

func (p *parser) parseFloat() (ast.Expression, error) {
	expr := &ast.FloatExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	val, err := strconv.ParseFloat(p.cur.Val, 64)
	if err != nil {
		return nil, p.errorf(p.cur, "float %s too large for a 64 bit float", p.cur.Val)
	}
	expr.Val = val
	return expr, nil
}

func (p *parser) parseBoolean() (ast.Expression, error) {
	return &ast.BooleanExpression{
		Val: p.cur.Val == "true",
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}, nil
}

func (p *parser) parseIdentifier() (ast.Expression, error) {
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
	}, nil
}

func (p *parser) parseCallExpression() (ast.Expression, error) {
	callExpr := &ast.CallExpression{
		Identifier: p.cur.Val,
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	if !p.expectPeek(lexer.LParen) {
		return nil, p.errorf(p.peek, "expected '(' received %s", p.peek.Val)
	}

	listErr := p.parseList(lexer.RParen, func() error {
		expr, err := p.parseExpression(lowest)
		if err != nil {
			return err
		}
		callExpr.Arguments = append(callExpr.Arguments, expr)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}
	return callExpr, nil
}

func (p *parser) parseIf() (ast.Expression, error) {
	// defer untrace(trace(fmt.Sprintf("IFTE %s", p.cur.Val)))

	expr := &ast.IfExpression{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	var err error
	// trace("ELSE")
	if expr.Condition, err = p.parseExpression(lowest); err != nil {
		return nil, err
	}
	// untrace("ELSE")

	if !p.expectPeek(lexer.Then) {
		return nil, p.errorf(p.peek, "expected 'then' received '%s'", p.peek.Val)
	}
	p.advance()

	if expr.Consequence, err = p.parseExpression(lowest); err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.Else) {
		return nil, p.errorf(p.peek, "expected 'else' received '%s'", p.peek.Val)
	}
	p.advance()

	if expr.Otherwise, err = p.parseExpression(lowest); err != nil {
		return nil, err
	}

	return expr, nil
}
