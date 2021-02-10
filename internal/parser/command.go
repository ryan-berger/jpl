package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

func (p *Parser) parseCommand() ast.Command {
	switch p.cur.Type {
	case lexer.Let:
		return p.parseLetStatement()
	case lexer.Function:
		return p.parseFunction()
	case lexer.Assert:
		return p.parseAssertStatement()
	default:
		return p.parseBuiltinCommand()
	}
}

func (p *Parser) parseAssertStatement() ast.Statement {
	stmt := &ast.AssertStatement{}
	p.advance()

	if stmt.Expr = p.parseExpression(lowest); stmt.Expr == nil {
		return nil
	}

	if !p.expectPeek(lexer.Comma) {
		return nil
	}

	if !p.expectPeek(lexer.String) {
		return nil
	}

	stmt.Message = p.cur.Val
	p.advance()
	return stmt
}

func (p *Parser) parseBuiltinCommand() ast.Command {
	if parse := p.cmdParseFns[p.cur.Type]; parse != nil {
		return parse()
	}
	p.errorf("err: expected command, received %s at line %d", p.cur.Val, p.cur.Line)
	return nil
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
	if !p.peekTokenIs(lexer.String) {
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

	return show
}

func (p *Parser) parseTimeCommand() ast.Command {
	time := &ast.Time{}
	p.advance()

	time.Command = p.parseCommand()
	if time.Command == nil {
		return nil
	}

	return time
}
