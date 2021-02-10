package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseFunction() ast.Command {
	function := &ast.Function{}

	if !p.expectPeek(lexer.Variable) {
		return nil
	}

	function.Var = p.cur.Val

	if !p.expectPeek(lexer.LParen) {
		return nil
	}

	function.Bindings = p.parseBindings()

	if !p.expectPeek(lexer.RParen) {
		return nil
	}

	if !p.expectPeek(lexer.Colon) {
		return nil
	}
	p.advance()

	if function.ReturnType = p.parseTypeExpression(); function.ReturnType == nil {
		return nil
	}

	if !p.expectPeek(lexer.LCurly) {
		return nil
	}

	if !p.expectPeek(lexer.NewLine) {
		return nil
	}
	p.advance()

	function.Statements = p.parseStatements()

	if !p.expectPeek(lexer.RCurly) {
		return nil
	}
	p.advance()
	return function
}

func (p *Parser) parseStatements() []ast.Statement {
	var statements []ast.Statement

	var stmt ast.Statement
	for !p.peekTokenIs(lexer.RCurly) && !p.peekTokenIs(lexer.EOF) {
		switch p.cur.Type {
		case lexer.Let:
			stmt = p.parseLetStatement()
		case lexer.Return:
			stmt = p.parseReturnStatement()
		case lexer.Assert: // TODO: actually implement
			stmt = p.parseAssertStatement()
		default:
			p.errorf("err :yeet") // TODO: YEET
			stmt = nil
			for !p.curTokenIs(lexer.NewLine) {
				p.advance()
			}

		}

		if stmt != nil {
			statements = append(statements, stmt)
		}

		p.advance()
	}

	return statements
}

func (p *Parser) parseBindings() []ast.Binding {
	var bindings []ast.Binding
	if p.expectPeek(lexer.RCurly) {
		p.advance() // move past rCurly
		return bindings
	}

	p.advance() // move past lCurly

	bindings = append(bindings, p.parseBinding()) // TODO: error handling
	for p.peekTokenIs(lexer.Comma) {
		p.advance()
		p.advance()
		bindings = append(bindings, p.parseBinding())
	}

	return bindings
}

func (p *Parser) parseTupleBinding() ast.Binding {
	return nil
}

func (p *Parser) parseBinding() ast.Binding {
	if p.curTokenIs(lexer.LCurly) {
		return p.parseTupleBinding()
	}

	binding := &ast.TypeBind{}
	binding.Argument = p.parseArgument()
	if binding.Argument == nil {
		return nil
	}

	if !p.expectPeek(lexer.Colon) {
		return nil
	}
	p.advance() // move past colon

	binding.Type = p.parseTypeExpression()
	if binding.Type == nil {
		return nil
	}

	return binding
}

var tokenToType = map[lexer.TokenType]ast.Type{
	lexer.Float: ast.Float,
	lexer.Int:   ast.Int,
}

