package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *parser) parseFunction() ast.Command {
	function := &ast.Function{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

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

	if !p.curTokenIs(lexer.RCurly) {
		return nil
	}
	p.advance()
	return function
}

func (p *parser) parseStatements() []ast.Statement {
	var statements []ast.Statement

	for !p.peekTokenIs(lexer.RCurly) && !p.peekTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt == nil {
			return nil
		}

		statements = append(statements, stmt)
		p.advance()
	}

	return statements
}

func (p *parser) parseBindings() []ast.Binding {
	var bindings []ast.Binding

	ok := p.parseList(lexer.RParen, func() bool {
		bind := p.parseBinding()
		if bind == nil {
			return false
		}
		bindings = append(bindings, bind)
		return true
	})

	if !ok {
		return nil
	}

	return bindings
}

func (p *parser) parseTupleBinding() ast.Binding {
	binding := &ast.TupleBinding{}

	ok := p.parseList(lexer.RCurly, func() bool {
		bind := p.parseBinding()
		if bind == nil {
			return false
		}
		binding.Bindings = append(binding.Bindings, bind)
		return true
	})

	if !ok {
		return nil
	}

	return binding
}

func (p *parser) parseBinding() ast.Binding {
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
	lexer.Bool:  ast.Boolean,
}
