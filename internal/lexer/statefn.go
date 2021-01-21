package lexer

import "unicode"

type stateFn func(l *lexer) stateFn

const (
	lCurly = "{"
	rCurly = "}"
	lParen = "("
	rParen = ")"
	lBrace = "["
	rBrace = "]"

	// stmt keywords
	function = "fn"
	let      = "let"
	assert   = "assert"

	// Operators
	equalTo            = "=="
	notEqualTo         = "!="
	lessThan           = "<"
	greaterThan        = ">"
	lessThanEqualTo    = "<="
	greaterThanEqualTo = ">="
	or                 = "||"
	and                = "&&"
	not                = "!"
	plus               = "+"
	minus              = "-"
	multiply           = "*"
	divide             = "/"
	mod                = "%"
)

//var keywords = map[string]TokenType {
//	"array",
//	"assert",
//	"bool",
//	"else",
//	"false",
//	"float3",
//	"float4",
//	"float",
//	"fn",
//	"if",
//	"int",
//	"let", "print", "read", "return", "show", "sum", "then", "time", "to", "true", "write"
//}

var infixOperators = []string{
	equalTo, notEqualTo, lessThan, greaterThan, lessThanEqualTo,
	greaterThanEqualTo, or, and, not, plus, minus, multiply, divide, mod,
}

// lexText is the default state that is returned at the end of lexing each state
func lexText(l *lexer) stateFn {
	for {
		switch next := l.next(); {
		case next == ' ':
			l.ignore()
		case !unicode.IsLetter(next):
			return l.errorf("reached invalid character: %s", string(next))
		case isAlphaNumeric(next):
			l.backup()
			return lexCmd
		}
	}
	return nil
}

func lexCmd(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			switch word {
			case let:
				l.emit(Let)
				return lexLet
			case function:
				l.emit(Function)
				return lexFunction
			case assert:
				l.emit(Assert)
				return lexAssert
			}
		}
	}
	return nil
}

func lexLet(l *lexer) stateFn      { return nil }
func lexFunction(l *lexer) stateFn { return nil }
func lexAssert(l *lexer) stateFn   { return nil }

func lexExpression(l *lexer) stateFn { return nil }
func lexLValue(l *lexer) stateFn     { return nil }

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
