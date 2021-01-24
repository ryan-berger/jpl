package lexer

import (
	"fmt"
	"strings"
)

// Lexer holds state while we lex. Heavily inspired by "Writing an Interpreter in Go"
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		ch:           0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isAlphabetic(l.ch) || isNumeric(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) errorf(msg string, args ...interface{}) Token {
	// jump to the end
	l.readPosition = len(l.input)
	l.ch = 0
	// return illegal token
	return newTokenString(ILLEGAL, fmt.Sprintf(msg, args...))
}

func (l *Lexer) readComment() *Token {
	l.readChar() // advance so we are in line with the comment
	for l.ch != '\n' {
		l.readChar()
	}
	l.readChar()
	return nil
}

func (l *Lexer) readMultilineComment() *Token {
	l.readChar() // advance so we are in line with the comment

	pos := l.position
	for !strings.HasSuffix(l.input[pos:l.position], "*/") {
		l.readChar()
	}

	l.readChar()
	return nil
}

func (l *Lexer) readString() Token {
	pos := l.position
	l.readChar() // advance past the first quotation mark

	for l.ch != '"' {
		if invalidChar(l.ch) || l.ch == '\n' {
			return l.errorf("error, expected end quote received: %s", string(l.ch))
		}
		l.readChar()
	}

	l.readChar()
	str := l.input[pos:l.position]
	return newTokenString(String, str)
}

func (l *Lexer) readDigits() string {
	pos := l.position
	if l.ch == '-' || l.ch == '+' {
		l.readChar() // move forward once if the first character is a m
	}
	for isNumeric(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() Token {
	first := l.readDigits()

	if l.ch != '.' {
		return newTokenString(IntLiteral, first)
	}

	// step ahead one character
	l.readChar()
	second := l.readDigits()
	return newTokenString(FloatLiteral, fmt.Sprintf("%s.%s", first, second))
}

func isAlphabetic(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isNumeric(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func invalidChar(ch byte) bool {
	return (ch < 32 || ch > 126) && ch != 10 && ch != 0
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func newToken(tokType TokenType, ch byte) Token {
	return Token{
		Type: tokType,
		Val:  string(ch),
	}
}

func newTokenString(tokenType TokenType, str string) Token {
	return Token{
		Type: tokenType,
		Val:  str,
	}
}

var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"if":     If,
	"then":   Then,
	"else":   Else,
	"return": Return,
	"array":  Array,
	// builtins
	"print":     Print,
	"show":      Show,
	"time":      Time,
	"sum":       Sum,
	"assert":    Assert,
	"read":      Read,
	"write":     Write,
	"to":        To,
	"attribute": Attribute,
}

func (l *Lexer) LexAll() ([]Token, bool) {
	tokens := make([]Token, 0)
	for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
		if len(tokens) != 0 && (tokens[len(tokens)-1].Type == NewLine && tok.Type == NewLine) {
			continue
		}
		tokens = append(tokens, tok)
	}
	tokens = append(tokens, Token{Type: EOF, Val: ""})
	return tokens, len(tokens) >= 2 && tokens[len(tokens)-2].Type != ILLEGAL
}

func (l *Lexer) NextToken() Token {
	var t Token

	// skip whitespace
	for l.ch == ' ' {
		l.readChar()
	}

	if (l.ch < 32 || l.ch > 126) && l.ch != 10 && l.ch != 0 {
		return l.errorf("invalid character: %s", string(l.ch))
	}

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			t = newTokenString(EqualTo, "==")
			l.readChar()
		} else {
			t = newToken(Assign, l.ch)
		}
	case '!':
		if l.peek() == '=' {
			t = newTokenString(NotEqualTo, "!=")
			l.readChar()
		} else {
			t = newToken(Not, l.ch)
		}
	case '>':
		if l.peek() == '=' {
			t = newTokenString(GreaterThanOrEqual, ">=")
			l.readChar()
		} else {
			t = newToken(GreaterThan, l.ch)
		}
	case '<':
		if l.peek() == '=' {
			t = newTokenString(LessThanOrEqual, "<=")
			l.readChar()
		} else {
			t = newToken(LessThan, l.ch)
		}
	case '&':
		if l.peek() == '&' {
			t = newTokenString(And, "&&")
			l.readChar()
		} else {
			t = newTokenString(ILLEGAL, fmt.Sprintf("expected & received %s", string(l.ch)))
		}
	case '|':
		if l.peek() == '|' {
			t = newTokenString(Or, "||")
			l.readChar()
		} else {
			t = newTokenString(ILLEGAL, fmt.Sprintf("expected | received %s", string(l.ch)))
		}
	case '+':
		if isNumeric(l.peek()) {
			return l.readNumber()
		}
		t = newToken(Plus, l.ch)
	case '-':
		if isNumeric(l.peek()) {
			return l.readNumber()
		}
		t = newToken(Minus, l.ch)
	case '*':
		t = newToken(Multiply, l.ch)
	case '/':
		peek := l.peek()
		if peek == '/' {
			if err := l.readComment(); err != nil {
				return *err
			}
			return l.NextToken()
		}
		if peek == '*' {
			if err := l.readMultilineComment(); err != nil {
				return *err
			}
			return l.NextToken()
		}
		t = newToken(Divide, l.ch)
	case '\\':
		if peek := l.peek(); peek != '\n' {
			return l.errorf("error, expected newline, received %s", string(peek))
		}
		l.readChar()
		l.readChar()
		return l.NextToken()
	case '%':
		t = newToken(Mod, l.ch)
	case '[':
		t = newToken(LBrace, l.ch)
	case ']':
		t = newToken(RBrace, l.ch)
	case '{':
		t = newToken(LCurly, l.ch)
	case '}':
		t = newToken(RCurly, l.ch)
	case '(':
		t = newToken(LParen, l.ch)
	case ')':
		t = newToken(RParen, l.ch)
	case ':':
		t = newToken(Colon, l.ch)
	case ',':
		t = newToken(Comma, l.ch)
	case '\n':
		l.readChar()
		for l.ch == ' ' || l.ch == '\n' {
			l.readChar()
		}
		return newToken(NewLine, '\n')
	case '"':
		return l.readString()
	case 0:
		t.Type = EOF
	default:
		if isAlphabetic(l.ch) {
			t.Type = Variable
			t.Val = l.readIdentifier()
			if tokenType, ok := keywords[t.Val]; ok {
				t.Type = tokenType
			}
			return t
		}
		if isNumeric(l.ch) {
			return l.readNumber()
		}
	}

	l.readChar()
	return t
}
