package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	let := &ast.LetStatement{}

	if let.LValue = p.parseLValue(); let.LValue == nil {
		return nil
	}

	if !p.expectPeek(lexer.Assign) {
		return nil
	}

	if let.Expr = p.parseExpression(); let.Expr == nil {
		return nil
	}

	return let
}

func (p *Parser) parseLValue() ast.LValue {
	switch {
	case p.curTokenIs(lexer.LCurly):
		p.advance()
		return p.parseTupleLValue()
	case p.curTokenIs(lexer.Variable):
		return p.parseArgument()
	}
	return nil
}

func (p *Parser) parseTupleLValue() ast.LValue {
	lTuple := &ast.LTuple{}
	for {
		var lVal ast.LValue
		switch {
		case p.curTokenIs(lexer.LCurly):
			p.advance()                 // advance token
			lVal = p.parseTupleLValue() // recurse to parse tuple val
		case p.curTokenIs(lexer.Variable):
			lVal = p.parseArgument() // parse arg
		default:
			return nil // if none of these are the case, return nil TODO: error out
		}

		if lVal == nil {
			return nil // return if there was an error in parsing either lVal
		}

		lTuple.Args = append(lTuple.Args, lVal) // append lVal
		if p.expectPeek(lexer.RBrace) {
			break // if we've found an RBrace, we are finished parsing
		}

		if !p.expectPeek(lexer.Comma) {
			return nil // if there isn't a comma, then the peek token is illegal
		}
	}

	return lTuple
}

func (p *Parser) parseArgument() ast.Argument {
	argName := p.cur.Val // TODO: check to make sure no keyword
	if !p.expectPeek(lexer.RBrace) {
		return &ast.VariableArgument{
			Variable: argName,
		}
	}
	return nil
}
