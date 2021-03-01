package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseLetStatement() ast.Statement {
	let := &ast.LetStatement{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	if let.LValue = p.parseLValue(); let.LValue == nil {
		return nil
	}

	if !p.expectPeek(lexer.Assign) {
		p.errorf("err: illegal token. Expected '=', found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	p.advance() // advance onto expression

	if let.Expr = p.parseExpression(lowest); let.Expr == nil { // get out of here if expression parsing fails
		return nil
	}

	p.advance()
	return let
}

func (p *Parser) parseLValue() ast.LValue {
	switch {
	case p.curTokenIs(lexer.LCurly):
		return p.parseTupleLValue()
	case p.curTokenIs(lexer.Variable):
		return p.parseArgument()
	}
	p.errorf("err: illegal token. Expected argument or '{', found %s at line %d", p.cur.Val, p.cur.Line)
	return nil
}

func (p *Parser) parseTupleLValue() ast.LValue {
	lTuple := &ast.LTuple{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	ok := p.parseList(lexer.RCurly, func() bool {
		expr := p.parseLValue()
		if expr == nil {
			return false
		}
		lTuple.Args = append(lTuple.Args, expr)
		return true
	})

	if !ok {
		return nil
	}

	return lTuple
}

func (p *Parser) parseArgument() ast.Argument {
	arg := &ast.VariableArr{
		Variable: p.cur.Val, // TODO: check to make sure no keyword
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	if !p.expectPeek(lexer.LBrace) {
		return &ast.VariableArgument{
			Variable: p.cur.Val,
		}
	}

	ok := p.parseList(lexer.RBrace, func() bool {
		if !p.curTokenIs(lexer.Variable) {
			return false
		}
		arg.Variables = append(arg.Variables, p.cur.Val)
		return true
	})

	if !ok {
		return nil
	}

	return arg
}
