package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Lex(text string) chan Token {
	return nil
}

const eof = -1

// lexer is a struct for internal state to help with the state machine
// and is heavily inspired by Rob Pike's talk on lexing in Go.
type lexer struct {
	name  string     // used only for error reports.
	input string     // the string being scanned.
	start int        // start position of the current item.
	pos   int        // pos is current position in the input.
	width int        // width of last rune read from input.
	items chan Token // channel of scanned items.
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t TokenType) {
	l.items <- Token{Type: t, Val: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Token{
		Type: Error,
		Val:  fmt.Sprintf(format, args...),
	}
	return nil
}

func (l *lexer) prefixed(with string) bool {
	return strings.HasPrefix(l.input[l.pos:], with)
}
