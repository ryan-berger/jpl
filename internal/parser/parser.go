package parser

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
	"github.com/ryan-berger/jpl/internal/errors"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/meta"
)

type parser struct {
	tokens   []lexer.Token
	position int
	cur      lexer.Token
	peek     lexer.Token

	error error

	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func newParser(tokens []lexer.Token) *parser {
	p := &parser{
		tokens:   tokens,
		position: 0,

		prefixParseFns: make(map[lexer.TokenType]prefixParseFn),
		infixParseFns:  make(map[lexer.TokenType]infixParseFn),
	}

	p.advance()
	p.advance()

	// "prefix" operators
	p.registerPrefixFn(lexer.Not, p.parsePrefixExpr)
	p.registerPrefixFn(lexer.LParen, p.parseGroupedExpression)
	p.registerPrefixFn(lexer.Minus, p.parsePrefixExpr)
	p.registerPrefixFn(lexer.LCurly, p.parseTupleExpression)
	p.registerPrefixFn(lexer.LBrace, p.parseArrayExpression)
	p.registerPrefixFn(lexer.IntLiteral, p.parseInteger)
	p.registerPrefixFn(lexer.FloatLiteral, p.parseFloat)
	p.registerPrefixFn(lexer.BoolLiteral, p.parseBoolean)
	p.registerPrefixFn(lexer.Variable, p.parseIdentifier)
	p.registerPrefixFn(lexer.If, p.parseIf)
	p.registerPrefixFn(lexer.Array, p.parseArrayTransform)
	p.registerPrefixFn(lexer.Sum, p.parseSumTransform)
	// type casts should be call expressions
	p.registerPrefixFn(lexer.Int, p.parseCallExpression)
	p.registerPrefixFn(lexer.Float, p.parseCallExpression)
	p.registerPrefixFn(lexer.Float3, p.parseCallExpression)
	p.registerPrefixFn(lexer.Float4, p.parseCallExpression)

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
	p.registerInfixFn(lexer.LessThanOrEqual, p.parseInfixExpr)
	p.registerInfixFn(lexer.GreaterThan, p.parseInfixExpr)
	p.registerInfixFn(lexer.GreaterThanOrEqual, p.parseInfixExpr)
	p.registerInfixFn(lexer.LCurly, p.parseTupleRefExpr)
	p.registerInfixFn(lexer.LBrace, p.parseArrayRefExpr)

	return p
}

func (p *parser) advance() {
	p.cur = p.peek
	p.peek = p.tokens[p.position]
	if p.position != len(p.tokens)-1 {
		p.position++
	}
}

func (p *parser) expectPeek(tokType lexer.TokenType) bool {
	if p.peek.Type == tokType {
		p.advance()
		return true
	}
	return false
}

func (p *parser) curTokenIs(tokType lexer.TokenType) bool {
	return p.cur.Type == tokType
}

func (p *parser) peekTokenIs(tokType lexer.TokenType) bool {
	return p.peek.Type == tokType
}

func (p *parser) registerPrefixFn(tokType lexer.TokenType, prefixFn prefixParseFn) {
	p.prefixParseFns[tokType] = prefixFn
}

func (p *parser) registerInfixFn(tokType lexer.TokenType, infixFn infixParseFn) {
	p.infixParseFns[tokType] = infixFn
}

func (p *parser) curPrecedence() precedence {
	if pr, ok := opPrecedence[p.cur.Type]; ok {
		return pr
	}
	return lowest
}
func (p *parser) peekPrecedence() precedence {
	if pr, ok := opPrecedence[p.peek.Type]; ok {
		return pr
	}
	return lowest
}

func (p *parser) errorf(loc meta.Locationer, format string, args ...interface{}) {
	p.error = errors.ParseError(fmt.Sprintf(format, args...), loc)
}

func (p *parser) parseProgram() ([]ast.Command, error) {
	var commands []ast.Command
	if p.curTokenIs(lexer.NewLine) {
		p.advance() // newline can be at the beginning of the file sometime
	}

	for !p.curTokenIs(lexer.EOF) {
		cmd := p.parseCommand()
		if cmd != nil {
			commands = append(commands, cmd)
		} else {
			return nil, p.error
		}

		p.advance()
	}

	return commands, nil
}

func Parse(tokens []lexer.Token) (ast.Program, error) {
	p := newParser(tokens)
	return p.parseProgram()
}