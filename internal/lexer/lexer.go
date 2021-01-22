package lexer

import "fmt"

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
	for isAlphabetic(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isAlphabetic(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
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
}

func (l *Lexer) NextToken() Token {
	var t Token

	// skip whitespace
	for l.ch == ' ' {
		l.readChar()
	}

	if (l.ch <= 32 || l.ch >= 126) && l.ch != 10 && l.ch != 0 {
		return newTokenString(ILLEGAL, fmt.Sprintf("invalid character: %s", string(l.ch)))
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
		t = newToken(Plus, l.ch)
	case '-':
		t = newToken(Minus, l.ch)
	case '*':
		t = newToken(Multiply, l.ch)
	case '/':
		t = newToken(Divide, l.ch)
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
		t = newToken(NewLine, l.ch)
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
	}

	l.readChar()
	return t
}
