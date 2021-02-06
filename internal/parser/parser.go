package parser

import (
	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

type Program struct {
	Commands []ast.Command
}

type Parser struct {
	tokens   []lexer.Token
	err      error
	position int
	cur      lexer.Token
	peek     lexer.Token
}

func NewParser(tokens []lexer.Token) *Parser {
	p := &Parser{
		tokens:   tokens,
		position: 0,
	}

	p.advance()
	p.advance()

	return p
}

func (p *Parser) advance() {
	p.cur = p.peek
	p.peek = p.tokens[p.position]
	p.position++
}

func (p *Parser) expectPeek(tokType lexer.TokenType) bool {
	if p.peek.Type == tokType {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) curTokenIs(tokType lexer.TokenType) bool {
	return p.cur.Type == tokType
}

func (p *Parser) ParseProgram() []ast.Command {
	var commands []ast.Command
	for p.cur.Type != lexer.EOF {
		cmd := p.parseCommand()
		if cmd != nil {
			commands = append(commands, cmd)
		}
		p.advance()
	}
	return commands
}

func (p *Parser) parseCommand() ast.Command {
	switch p.cur.Type {
	case lexer.Let:
		return p.parseLetStatement()
	case
	default:
		return nil
	}
}
