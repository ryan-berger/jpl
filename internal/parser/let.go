package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseLetStatement() ast.Statement {
	let := &ast.LetStatement{}
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
	lTuple := &ast.LTuple{}

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
	argName := p.cur.Val // TODO: check to make sure no keyword
	if !p.expectPeek(lexer.LBrace) {
		return &ast.VariableArgument{
			Variable: argName,
		}
	}

	var args []string
	ok := p.parseList(lexer.RBrace, func() bool {
		if !p.curTokenIs(lexer.Variable) {
			return false
		}
		args = append(args, p.cur.Val)
		return true
	})

	if !ok {
		return nil
	}

	return &ast.VariableArr{
		Variable:  argName,
		Variables: args,
	}
}
