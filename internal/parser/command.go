package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *parser) parseCommand() (ast.Command, error) {
	switch p.cur.Type {
	case lexer.Read:
		return p.parseReadCommand()
	case lexer.Write:
		return p.parseWriteCommand()
	case lexer.Print:
		return p.parsePrintCommand()
	case lexer.Show:
		return p.parseShowCommand()
	case lexer.Time:
		return p.parseTimeCommand()
	case lexer.Function:
		return p.parseFunction()
	default:
		return p.parseStatement()
	}
}

func (p *parser) parseStatement() (ast.Statement, error) {
	switch p.cur.Type {
	case lexer.Let:
		return p.parseLetStatement()
	case lexer.Return:
		return p.parseReturnStatement()
	case lexer.Assert:
		return p.parseAssertStatement()
	default:
		return nil, p.errorf(p.cur, "expected let, return, assert, received: %s", p.cur.Val)
	}
}

func (p *parser) parseAssertStatement() (ast.Statement, error) {
	stmt := &ast.AssertStatement{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	var err error
	if stmt.Expr, err = p.parseExpression(lowest); err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.Comma) {
		return nil, p.errorf(p.peek, "expected comma received: %s", p.peek.Val)
	}

	if !p.expectPeek(lexer.String) {
		return nil, p.errorf(p.peek, "expected string received: %s", p.peek.Val)
	}

	stmt.Message = p.cur.Val
	p.advance()

	return stmt, nil
}

func (p *parser) parseReadCommand() (ast.Command, error) {
	read := &ast.Read{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	if !p.expectPeek(lexer.Variable) {
		return nil, p.errorf(p.peek, "illegal token expected read type found %s", p.peek.Val)
	}

	if p.cur.Val != "image" && p.cur.Val != "video" {
		return nil, p.errorf(p.peek, "err: unsupported read type %s", p.peek.Val)
	}

	read.Type = p.cur.Val

	if !p.expectPeek(lexer.String) {
		return nil, p.errorf(p.peek, "illegal token. Expected string, found %s", p.peek.Val)
	}

	var err error
	read.Src, err = p.parseExpression(lowest)
	if err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.To) {
		return nil, p.errorf(p.peek, "err: illegal token. Expected 'to', found %s", p.peek.Val)
	}
	p.advance()

	read.Argument, err = p.parseArgument()
	if read.Argument == nil {
		return nil, err
	}

	p.advance()
	return read, nil
}

func (p *parser) parseWriteCommand() (ast.Command, error) {
	write := &ast.Write{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	if !p.expectPeek(lexer.Variable) {
		return nil, p.errorf(p.peek, "illegal token. Expected write type, found %s", p.peek.Val)
	}

	if p.cur.Val != "image" && p.cur.Val != "video" {
		return nil, p.errorf(p.peek, "err: unsupported write type %s", p.peek.Val)
	}

	write.Type = p.cur.Val
	p.advance()

	var err error
	write.Expr, err = p.parseExpression(lowest)
	if err != nil {
		return nil, err
	}

	if !p.expectPeek(lexer.To) {
		return nil, p.errorf(p.peek, "illegal token. Expected 'to', found %s", p.peek.Val)
	}
	p.advance()

	write.Dest, err = p.parseExpression(lowest)
	if err != nil {
		return nil, err
	}

	p.advance()
	return write, nil
}

func (p *parser) parsePrintCommand() (ast.Command, error) {
	pr := &ast.Print{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	if !p.expectPeek(lexer.String) {
		return nil, p.errorf(p.peek, "illegal token. Expected string, found %s", p.peek.Val)
	}

	pr.Str = p.cur.Val
	p.advance()
	return pr, nil
}

func (p *parser) parseShowCommand() (ast.Command, error) {
	show := &ast.Show{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	var err error
	show.Expr, err = p.parseExpression(lowest)
	if err != nil {
		return nil, err
	}

	p.advance()
	return show, nil
}

func (p *parser) parseTimeCommand() (ast.Command, error) {
	time := &ast.Time{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	var err error
	time.Command, err = p.parseCommand()
	if err != nil {
		return nil, err
	}

	if !p.curTokenIs(lexer.NewLine) {
		p.advance()
	}

	return time, nil
}
