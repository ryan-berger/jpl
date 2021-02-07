package parser

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/lexer"
)

type Program struct {
	Commands []ast.Command
}

type Parser struct {
	tokens   []lexer.Token
	position int
	cur      lexer.Token
	peek     lexer.Token

	errors []error

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn

	cmdParseFns map[lexer.TokenType]func() ast.Command
}

func NewParser(tokens []lexer.Token) *Parser {
	p := &Parser{
		tokens:   tokens,
		position: 0,

		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),

		cmdParseFns: make(map[lexer.TokenType]func() ast.Command),
	}

	p.advance()
	p.advance()

	// "prefix" operators
	p.registerPrefixFn(lexer.Not, p.parsePrefixExpr)
	p.registerPrefixFn(lexer.LParen, p.parseGroupedExpression)
	p.registerPrefixFn(lexer.Minus, p.parsePrefixExpr)
	p.registerPrefixFn(lexer.LCurly, p.parseTupleExpression)
	p.registerPrefixFn(lexer.IntLiteral, p.parseInteger)
	p.registerPrefixFn(lexer.FloatLiteral, p.parseFloat)
	p.registerPrefixFn(lexer.Bool, p.parseBoolean)
	p.registerPrefixFn(lexer.Variable, p.parseIdentifier)
	p.registerPrefixFn(lexer.If, p.parseIf)
	// p.registerPrefixFn(lexer.Array, p.parseArr)
	// p.registerPrefixFn(lexer.Sum, p.parseSum)

	p.registerInfixFn(lexer.Plus, p.parseInfixExpr)
	p.registerInfixFn(lexer.Minus, p.parseInfixExpr)
	p.registerInfixFn(lexer.Multiply, p.parseInfixExpr)
	p.registerInfixFn(lexer.Divide, p.parseInfixExpr)
	p.registerInfixFn(lexer.Mod, p.parseInfixExpr)
	p.registerInfixFn(lexer.And, p.parseInfixExpr)
	p.registerInfixFn(lexer.Or, p.parseInfixExpr)
	p.registerInfixFn(lexer.EqualTo, p.parseInfixExpr)
	p.registerInfixFn(lexer.NotEqualTo, p.parseInfixExpr)
	p.registerInfixFn(lexer.LessThan, p.parseInfixExpr)
	p.registerInfixFn(lexer.GreaterThan, p.parseInfixExpr)
	p.registerInfixFn(lexer.GreaterThanOrEqual, p.parseInfixExpr)
	p.registerInfixFn(lexer.LessThanOrEqual, p.parseInfixExpr)
	// p.registerInfixFn(lexer.LBrace, p.parseInfixExpr)
	// p.registerInfixFn(lexer.LCurly, p.parseInfixExpr)

	p.registerCommandFn(lexer.Read, p.parseReadCommand)
	p.registerCommandFn(lexer.Write, p.parseWriteCommand)
	p.registerCommandFn(lexer.Print, p.parsePrintCommand)
	p.registerCommandFn(lexer.Show, p.parseShowCommand)
	p.registerCommandFn(lexer.Time, p.parseTimeCommand)

	return p
}

func (p *Parser) advance() {
	p.cur = p.peek
	p.peek = p.tokens[p.position]
	if p.position != len(p.tokens)-1 {
		p.position++
	}
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

func (p *Parser) peekTokenIs(tokType lexer.TokenType) bool {
	return p.peek.Type == tokType
}

func (p *Parser) registerPrefixFn(tokType lexer.TokenType, prefixFn prefixParseFn) {
	p.prefixParseFns[tokType] = prefixFn
}

func (p *Parser) registerInfixFn(tokType lexer.TokenType, infixFn infixParseFn) {
	p.infixParseFns[tokType] = infixFn
}

func (p *Parser) registerCommandFn(cmd lexer.TokenType, parseFn func() ast.Command) {
	p.cmdParseFns[cmd] = parseFn
}

func (p *Parser) curPrecedence() precedence {
	if pr, ok := opPrecedence[p.cur.Type]; ok {
		return pr
	}
	return lowest
}
func (p *Parser) peekPrecedence() precedence {
	if pr, ok := opPrecedence[p.peek.Type]; ok {
		return pr
	}
	return lowest
}

func (p *Parser) errorf(format string, args ...interface{}) {
	p.errors = append(p.errors, fmt.Errorf(format, args...))
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
	for _, cmd := range commands {
		fmt.Print(cmd.String())
	}

	return commands
}
