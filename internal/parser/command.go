package parser

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *parser) parseCommand() ast.Command {
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
		stmt := p.parseStatement()
		if stmt != nil {
			return stmt
		}
		p.error = fmt.Errorf("error while parsing command statement %w", p.error)
		return nil
	}
}

func (p *parser) parseStatement() ast.Statement {
	switch p.cur.Type {
	case lexer.Let:
		return p.parseLetStatement()
	case lexer.Return:
		return p.parseReturnStatement()
	case lexer.Assert:
		return p.parseAssertStatement()
	default:
		return nil
	}
}

func (p *parser) parseAssertStatement() ast.Statement {
	stmt := &ast.AssertStatement{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	if stmt.Expr = p.parseExpression(lowest); stmt.Expr == nil {
		return nil
	}

	if !p.expectPeek(lexer.Comma) {
		p.errorf(p.peek, "expected comma received: %s", p.peek.Val)
		return nil
	}

	if !p.expectPeek(lexer.String) {
		p.errorf(p.peek, "expected string received: %s", p.peek.Val)
		return nil
	}

	stmt.Message = p.cur.Val
	p.advance()
	return stmt
}

func (p *parser) parseReadCommand() ast.Command {
	read := &ast.Read{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	if !p.expectPeek(lexer.Variable) {
		p.errorf(p.peek, "illegal token expected read type found %s", p.peek.Val)
		return nil
	}

	if p.cur.Val != "image" && p.cur.Val != "video" {
		p.errorf(p.peek, "err: unsupported read type %s", p.peek.Val)
		return nil
	}

	read.Type = p.cur.Val

	if !p.expectPeek(lexer.String) {
		p.errorf(p.peek, "illegal token. Expected string, found %s", p.peek.Val)
		return nil
	}

	read.Src = p.cur.Val

	if !p.expectPeek(lexer.To) {
		p.errorf(p.peek,"err: illegal token. Expected 'to', found %s", p.peek.Val)
		return nil
	}
	p.advance()

	read.Argument = p.parseArgument()
	if read.Argument == nil {
		return nil
	}

	p.advance()
	return read
}

func (p *parser) parseWriteCommand() ast.Command {
	write := &ast.Write{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	if !p.expectPeek(lexer.Variable) {
		p.errorf(p.peek,"illegal token. Expected write type, found %s", p.peek.Val)
		return nil
	}

	if p.cur.Val != "image" && p.cur.Val != "video" {
		p.errorf(p.peek,"err: unsupported write type %s", p.peek.Val)
		return nil
	}

	write.Type = p.cur.Val
	p.advance()

	write.Expr = p.parseExpression(lowest)
	if write.Expr == nil {
		return nil
	}

	if !p.expectPeek(lexer.To) {
		p.errorf(p.peek, "illegal token. Expected 'to', found %s", p.peek.Val)
		return nil
	}

	if !p.expectPeek(lexer.String) {
		p.errorf(p.peek, "illegal token. Expected string, found %s", p.peek.Val)
		return nil
	}
	write.Dest = p.cur.Val

	p.advance()
	return write
}

func (p *parser) parsePrintCommand() ast.Command {
	pr := &ast.Print{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}

	if !p.expectPeek(lexer.String) {
		p.errorf(p.peek, "illegal token. Expected string, found %s", p.peek.Val)
		return nil
	}

	pr.Str = p.cur.Val
	p.advance()
	return pr
}

func (p *parser) parseShowCommand() ast.Command {
	show := &ast.Show{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	show.Expr = p.parseExpression(lowest)
	if show.Expr == nil {
		return nil
	}

	p.advance()
	return show
}

func (p *parser) parseTimeCommand() ast.Command {
	time := &ast.Time{
		Location: ast.Location{
			Line: p.cur.Line,
			Pos:  p.cur.Character,
		},
	}
	p.advance()

	time.Command = p.parseCommand()
	if time.Command == nil {
		return nil
	}

	if !p.curTokenIs(lexer.NewLine) {
		p.advance()
	}

	return time
}
