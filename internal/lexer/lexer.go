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
	lineNumber   int
	linePos      int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		lineNumber:   1,
		linePos:      0,
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
	l.linePos++
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isAlphabetic(l.ch) || isNumeric(l.ch) || l.ch == '_' || l.ch == '.' {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) errorf(msg string, args ...interface{}) Token {
	// jump to the end
	l.readPosition = len(l.input) + 1
	l.ch = 0
	// return illegal token
	return Token{
		Type:      ILLEGAL,
		Val:       fmt.Sprintf("%s at line %d position %d", fmt.Sprintf(msg, args...), l.lineNumber, l.linePos),
		Line:      l.lineNumber,
		Character: l.linePos,
	}
}

// readComment reads a single line comment
func (l *Lexer) readComment() {
	l.readChar() // advance so we are in line with the comment
	for l.ch != '\n' && l.readPosition <= len(l.input) {
		l.readChar()
	}
}

// readMultilineComment reads a multi-line comment
func (l *Lexer) readMultilineComment() *Token {
	l.readChar()
	l.readChar() // advance so we are in line with the comment

	pos := l.position
	for !strings.HasSuffix(l.input[pos:l.position], "*/") && l.readPosition <= len(l.input) {
		if invalidChar(l.ch) {
			tok := l.errorf("error, expected closing '*/' received invalid character: '%s'", string(l.ch))
			return &tok
		}
		if l.ch == '\n' {
			l.lineNumber++
			l.linePos = 0
		}
		l.readChar()
	}

	if !strings.HasSuffix(l.input[pos:l.position], "*/") {
		tok := l.errorf("error, expected closing '*/' received EOF")
		return &tok
	}

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
	return l.newOverflowedToken(String, str)
}

func (l *Lexer) readDigits() string {
	pos := l.position
	for isNumeric(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() Token {
	first := l.readDigits()

	if l.ch != '.' {
		return l.newOverflowedToken(IntLiteral, first)
	}

	// step ahead one character
	l.readChar()
	second := l.readDigits()
	return l.newOverflowedToken(FloatLiteral, fmt.Sprintf("%s.%s", first, second))
}

func isAlphabetic(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumeric(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func invalidChar(ch byte) bool {
	return (ch < 32 || ch > 126) && ch != 10
}

func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) newToken(tokType TokenType, ch byte) Token {
	return Token{
		Type:      tokType,
		Val:       string(ch),
		Line:      l.lineNumber,
		Character: l.linePos - 1,
	}
}

func (l *Lexer) newOverflowedToken(tokType TokenType, str string) Token {
	return Token{
		Type:      tokType,
		Val:       str,
		Line:      l.lineNumber,
		Character: l.linePos - len(str) - 1,
	}
}

func (l *Lexer) newTokenString(tokType TokenType, str string) Token {
	return Token{
		Type:      tokType,
		Val:       str,
		Line:      l.lineNumber,
		Character: l.linePos - len(str),
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
	// types
	"int":    Int,
	"float":  Float,
	"float3": Float3,
	"float4": Float4,
}

func (l *Lexer) LexAll() ([]Token, bool) {
	tokens := make([]Token, 0)
	for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
		if len(tokens) != 0 && (tokens[len(tokens)-1].Type == NewLine && tok.Type == NewLine) {
			continue
		}
		tokens = append(tokens, tok)
	}
	// TODO: Don't do this, just whyyyyy
	tokens = append(tokens, Token{Type: EOF, Val: ""})
	return tokens, len(tokens) == 1 || len(tokens) >= 2 && tokens[len(tokens)-2].Type != ILLEGAL
}

// searchNextToken looks for the next token for the given input
// this allows the lexer to be much more specific and avoid recursion.
// all characters such as /, /*, and \ that can be ignored are filtered
// out in order to avoid unnecessary recursion
// if a lexical error is encountered, an error token is returned for use
// in the NextToken function
func (l *Lexer) searchNextToken() *Token {
	for {
		switch l.ch {
		case ' ':
			l.readChar()
		case '/':
			if l.peek() == '/' {
				l.readComment()
			} else if l.peek() == '*' {
				if tok := l.readMultilineComment(); tok != nil {
					return tok
				}
			} else {
				return nil
			}
		case '\\':
			if l.peek() != '\n' {
				tok := l.errorf("error, expected newline after '\\' received: '%s'",
					string(l.peek()))
				return &tok
			}
			l.readChar() // advance to newline character
			l.readChar() // advance past newline character
			l.lineNumber++
			l.linePos = 0
		default: // we don't have anything to skip
			return nil
		}
	}
}

func (l *Lexer) NextToken() Token {
	var t Token

	// make sure we are in line with the next token
	if tok := l.searchNextToken(); tok != nil {
		return *tok
	}

	// if we've past the amount of input we have, we are at an EOF
	if l.readPosition > len(l.input) {
		return Token{Type: EOF}
	}

	// if we've reached an invalid character, we need to return
	if invalidChar(l.ch) {
		return l.errorf("error, received invalid character: %s", string(l.ch))
	}

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			l.readChar()
			t = l.newTokenString(EqualTo, "==")
		} else {
			t = l.newToken(Assign, l.ch)
		}
	case '!':
		if l.peek() == '=' {
			l.readChar()
			t = l.newTokenString(NotEqualTo, "!=")
		} else {
			t = l.newToken(Not, l.ch)
		}
	case '>':
		if l.peek() == '=' {
			l.readChar()
			t = l.newTokenString(GreaterThanOrEqual, ">=")
		} else {
			t = l.newToken(GreaterThan, l.ch)
		}
	case '<':
		if l.peek() == '=' {
			l.readChar()
			t = l.newTokenString(LessThanOrEqual, "<=")
		} else {
			t = l.newToken(LessThan, l.ch)
		}
	case '&':
		if l.peek() == '&' {
			l.readChar()
			t = l.newTokenString(And, "&&")
		} else {
			return l.errorf("error, expected '&' received '%s'", string(l.peek()))
		}
	case '|':
		if l.peek() == '|' {
			l.readChar()
			t = l.newTokenString(Or, "||")
		} else {
			return l.errorf("error, expected '|' received '%s'", string(l.peek()))
		}
	case '.':
		if !isNumeric(l.peek()) {
			return l.errorf("error, expected number received %s", string(l.peek()))
		}
		l.readChar()
		return l.newTokenString(FloatLiteral, fmt.Sprintf(".%s", l.readDigits()))
	case '+':
		t = l.newToken(Plus, l.ch)
	case '-':
		t = l.newToken(Minus, l.ch)
	case '*':
		t = l.newToken(Multiply, l.ch)
	case '/':
		t = l.newToken(Divide, l.ch)
	case '%':
		t = l.newToken(Mod, l.ch)
	case '[':
		t = l.newToken(LBrace, l.ch)
	case ']':
		t = l.newToken(RBrace, l.ch)
	case '{':
		t = l.newToken(LCurly, l.ch)
	case '}':
		t = l.newToken(RCurly, l.ch)
	case '(':
		t = l.newToken(LParen, l.ch)
	case ')':
		t = l.newToken(RParen, l.ch)
	case ':':
		t = l.newToken(Colon, l.ch)
	case ',':
		t = l.newToken(Comma, l.ch)
	case '\n':
		t = l.newToken(NewLine, '\n')
		l.lineNumber++
		l.linePos = 0
	case '"':
		return l.readString()
	default:
		if isAlphabetic(l.ch) {
			typ := Variable
			val := l.readIdentifier()
			if tokenType, ok := keywords[val]; ok {
				typ = tokenType
			}
			return l.newOverflowedToken(typ, val)
		}
		if isNumeric(l.ch) {
			return l.readNumber()
		}

		l.errorf("error, no token match to token: %s", string(l.ch))
	}

	l.readChar()
	return t
}
