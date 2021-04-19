package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/types"
)

func (p *parser) parseFunction() (ast.Command, error) {
	function := &ast.Function{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	if !p.expectPeek(lexer.Variable) {
		return nil, p.errorf(p.peek, "expected identifier, received %s", p.peek.Val)
	}

	function.Var = p.cur.Val

	if !p.expectPeek(lexer.LParen) {
		return nil, p.errorf(p.peek, "expected '(' received %s", p.peek.Val)
	}

	var err error
	function.Bindings, err = p.parseBindings()

	if err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.Colon) {
		return nil, p.errorf(p.peek, "expected ':' received %s", p.peek.Val)
	}
	p.advance()


	if function.ReturnType, err = p.parseTypeExpression(); function.ReturnType == nil {
		return nil, err
	}

	if !p.expectPeek(lexer.LCurly) {
		return nil, p.errorf(p.peek, "expected '{', received %s", p.peek.Val)
	}

	if !p.expectPeek(lexer.NewLine) {
		return nil, p.errorf(p.peek, "expected newline, received %s", p.peek.Val)
	}
	p.advance()

	stmts, err := p.parseStatements()
	if err != nil {
		return nil, err
	}

	function.Statements = stmts
	p.advance()
	return function, nil
}

func (p *parser) parseStatements() ([]ast.Statement, error) {
	var statements []ast.Statement

	for !p.curTokenIs(lexer.RCurly) && !p.curTokenIs(lexer.EOF) {
		stmt, err := p.parseStatement()

		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
		p.advance()
	}

	return statements, nil
}

func (p *parser) parseBindings() ([]ast.Binding, error) {
	var bindings []ast.Binding

	listErr := p.parseList(lexer.RParen, func() error {
		bind, err := p.parseBinding()
		if err != nil {
			return err
		}
		bindings = append(bindings, bind)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}

	return bindings, nil
}

func (p *parser) parseTupleBinding() (ast.Binding, error) {
	binding := &ast.TupleBinding{}

	listErr := p.parseList(lexer.RCurly, func() error {
		bind, err := p.parseBinding()
		if err != nil {
			return err
		}
		binding.Bindings = append(binding.Bindings, bind)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}

	return binding, nil
}

func (p *parser) parseBinding() (ast.Binding, error) {
	if p.curTokenIs(lexer.LCurly) {
		return p.parseTupleBinding()
	}

	var err error
	binding := &ast.TypeBind{}
	if binding.Argument, err = p.parseArgument(); err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.Colon) {
		return nil, p.errorf(p.peek, "expected ':' received %s", p.peek.Val)
	}
	p.advance() // move past colon

	if binding.Type, err = p.parseTypeExpression(); err != nil {
		return nil, err
	}

	return binding, nil
}

var tokenToType = map[lexer.TokenType]types.Type{
	lexer.Float: types.Float,
	lexer.Int:   types.Integer,
	lexer.Bool:  types.Boolean,
}
