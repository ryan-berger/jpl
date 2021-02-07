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

	if !p.expectPeek(lexer.Colon) {
		return nil
	}
	p.advance()

	if function.ReturnType = p.parseType(); function.ReturnType == nil {
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
			break
		}
		statements = append(statements, stmt)
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

	if !p.expectPeek(lexer.RParen) {
		return nil
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

	binding.Type = p.parseType()
	if binding.Type == nil {
		return nil
	}

	return binding
}

var tokenToType = map[lexer.TokenType]ast.Type{
	lexer.Float: ast.Float,
	lexer.Int:   ast.Int,
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

	if !p.expectPeek(lexer.LBrace) {
		return t
	}

	rank := 1
	for p.peekTokenIs(lexer.Comma) {
		rank++
	}

	if !p.expectPeek(lexer.RBrace) {
		return nil
	}

	return &ast.ArrType{
		Type: t,
		Rank: rank,
	}
}

func (p *Parser) parseTupleType() ast.Type {
	tupleType := &ast.TupleType{}

	if p.expectPeek(lexer.RCurly) {
		p.advance() // move past rCurly
		return nil
	}

	p.advance() // move past lCurly

	tupleType.Types = append(tupleType.Types, p.parseType()) // TODO: error handling
	for p.peekTokenIs(lexer.Comma) {
		p.advance()
		p.advance()
		tupleType.Types = append(tupleType.Types, p.parseType()) // TODO: error handling
	}

	return tupleType
}
