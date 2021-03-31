package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/types"
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
		p.error = NewError(p.peek, "expected ':' received %s", p.peek.Val)
		return nil
	}
	p.advance()

	if function.ReturnType = p.parseTypeExpression(); function.ReturnType == nil {
		return nil
	}

	if !p.expectPeek(lexer.LCurly) {
		p.errorf(p.peek, "expected '{', received %s", p.peek.Val)
		return nil
	}

	if !p.expectPeek(lexer.NewLine) {
		p.errorf(p.peek, "expected newline, received %s", p.peek.Val)
		return nil
	}
	p.advance()

	stmts, err := p.parseStatements()
	if err != nil {
		return nil
	}

	function.Statements = stmts
	p.advance()
	return function
}

func (p *parser) parseStatements() ([]ast.Statement, error) {
	var statements []ast.Statement

	for !p.curTokenIs(lexer.RCurly) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt == nil {
			if p.error == nil {
				panic("error not handled somewhere")
			}
			return nil, nil
		}

		statements = append(statements, stmt)
		p.advance()
	}

	return statements, nil
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

var tokenToType = map[lexer.TokenType]types.Type{
	lexer.Float: types.Float,
	lexer.Int:   types.Integer,
	lexer.Bool:  types.Boolean,
}
