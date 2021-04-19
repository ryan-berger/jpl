package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *parser) parseLetStatement() (ast.Statement, error) {
	// defer untrace(trace("LET"))
	let := &ast.LetStatement{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	var err error
	if let.LValue, err = p.parseLValue(); let.LValue == nil {
		return nil, err
	}

	if !p.expectPeek(lexer.Assign) {
		return nil, p.errorf(p.peek, "illegal token. Expected '=', found %s", p.peek.Val)
	}

	p.advance() // advance onto expression

	if let.Expr, err = p.parseExpression(lowest); err != nil { // get out of here if expression parsing fails
		return nil, err
	}

	p.advance()
	return let, nil
}

func (p *parser) parseLValue() (ast.LValue, error) {
	switch {
	case p.curTokenIs(lexer.LCurly):
		return p.parseTupleLValue()
	case p.curTokenIs(lexer.Variable):
		return p.parseArgument()
	}
	return nil, p.errorf(p.cur, "illegal token. Expected argument or '{', found '%s'", p.cur.Val)
}

func (p *parser) parseTupleLValue() (ast.LValue, error) {
	lTuple := &ast.LTuple{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	listErr := p.parseList(lexer.RCurly, func() error {
		expr, err := p.parseLValue()
		if err != nil {
			return err
		}
		lTuple.Args = append(lTuple.Args, expr)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}

	return lTuple, nil
}

func (p *parser) parseArgument() (ast.Argument, error) {
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
		}, nil
	}

	err := p.parseList(lexer.RBrace, func() error {
		if !p.curTokenIs(lexer.Variable) {
			return p.errorf(p.cur,
				"err, expected variable received %s", p.cur.Val)
		}

		arg.Variables = append(arg.Variables, p.cur.Val)
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(arg.Variables) == 0 {
		return nil, p.errorf(arg, "expected identifiers, received ]")
	}
	return arg, nil
}
