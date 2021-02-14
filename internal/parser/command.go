package parser

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseCommand() ast.Command {
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
	default:
		stmt := p.parseStatement()
		if stmt != nil {
			return stmt
		}
		p.error = fmt.Errorf("error while parsing command statement %w", p.error)
		return nil
	}
}

func (p *Parser) parseStatement() ast.Statement {
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

func (p *Parser) parseAssertStatement() ast.Statement {
	stmt := &ast.AssertStatement{}
	p.advance()

	if stmt.Expr = p.parseExpression(lowest); stmt.Expr == nil {
		return nil
	}

	if !p.expectPeek(lexer.Comma) {
		p.errorf("error, expected comma received: %s, line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	if !p.expectPeek(lexer.String) {
		p.errorf("error, expected string received: %s, line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	stmt.Message = p.cur.Val
	p.advance()
	return stmt
}

func (p *Parser) parseReadCommand() ast.Command {
	read := &ast.Read{}
	if !p.expectPeek(lexer.Variable) {
		p.errorf("err: illegal token. Expected read type, found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	if p.cur.Val != "image" && p.cur.Val != "video" {
		p.errorf("err: unsupported read type %s, at %d", p.peek.Val, p.peek.Line)
		return nil
	}

	read.Type = p.cur.Val

	if !p.expectPeek(lexer.String) {
		p.errorf("err: illegal token. Expected string, found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	read.Src = p.cur.Val

	if !p.expectPeek(lexer.To) {
		p.errorf("err: illegal token. Expected 'to', found %s at line %d", p.peek.Val, p.peek.Line)
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

func (p *Parser) parseWriteCommand() ast.Command {
	write := &ast.Write{}
	if !p.expectPeek(lexer.Variable) {
		p.errorf("err: illegal token. Expected write type, found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	if p.cur.Val != "image" && p.cur.Val != "video" {
		p.errorf("err: unsupported write type %s, at %d", p.peek.Val, p.peek.Line)
		return nil
	}

	write.Type = p.cur.Val
	p.advance()

	write.Expr = p.parseExpression(lowest)
	if write.Expr == nil {
		return nil
	}

	if !p.expectPeek(lexer.To) {
		p.errorf("err: illegal token. Expected 'to', found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	if !p.expectPeek(lexer.String) {
		p.errorf("err: illegal token. Expected string, found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}
	write.Dest = p.cur.Val

	p.advance()
	return write
}

func (p *Parser) parsePrintCommand() ast.Command {
	pr := &ast.Print{}
	if !p.expectPeek(lexer.String) {
		p.errorf("err: illegal token. Expected string, found %s at line %d", p.peek.Val, p.peek.Line)
		return nil
	}

	pr.Str = p.cur.Val
	p.advance()
	return pr
}

func (p *Parser) parseShowCommand() ast.Command {
	show := &ast.Show{}
	p.advance()

	show.Expr = p.parseExpression(lowest)
	if show.Expr == nil {
		return nil
	}

	p.advance()
	return show
}

func (p *Parser) parseTimeCommand() ast.Command {
	time := &ast.Time{}
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
