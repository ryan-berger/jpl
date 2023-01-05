package lexer

import "fmt"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Delimiters
	LParen
	RParen
	LBrace
	RBrace
	LCurly
	RCurly
	Comma
	Colon
	NewLine

	// Operators
	Assign
	Plus
	Minus
	Divide
	Multiply
	Mod
	Not
	And
	Or
	LessThan
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
	EqualTo
	NotEqualTo

	// Keywords
	Function
	Let
	If
	Then
	Else
	Return

	Variable

	// Builtins
	Read
	Write
	Assert
	Array
	Show
	Time
	Attribute
	To
	Sum
	Print

	// Literals
	String
	IntLiteral
	FloatLiteral
	BoolLiteral

	// Types
	Int
	Float
	Float3
	Float4
	Bool
	Pict
	Str
)

type Token struct {
	Type      TokenType
	Val       string
	Line      int
	Character int
}

func (t Token) GoString() string {
	return fmt.Sprintf("Token{Type: %s, Val: %s, Line: %d, Character: %d}",
		typeToDump[t.Type], t.Val, t.Line, t.Character)
}

func (t Token) Loc() (int, int) {
	return t.Line, t.Character
}

var typeToDump = map[TokenType]string{
	ILLEGAL: "ERROR",
	EOF:     "END_OF_FILE",

	// Delimiters
	LParen:  "LPAREN",
	RParen:  "RPAREN",
	LBrace:  "LSQUARE",
	RBrace:  "RSQUARE",
	LCurly:  "LCURLY",
	RCurly:  "RCURLY",
	Comma:   "COMMA",
	Colon:   "COLON",
	NewLine: "NEWLINE",

	// Operators
	Assign:             "EQUALS",
	Plus:               "OP",
	Minus:              "OP",
	Divide:             "OP",
	Multiply:           "OP",
	Mod:                "OP",
	Not:                "OP",
	And:                "OP",
	Or:                 "OP",
	LessThan:           "OP",
	LessThanOrEqual:    "OP",
	GreaterThan:        "OP",
	GreaterThanOrEqual: "OP",
	EqualTo:            "OP",
	NotEqualTo:         "OP",

	// Keywords
	Function: "FN",
	Let:      "LET",
	If:       "IF",
	Then:     "THEN",
	Else:     "ELSE",
	Return:   "RETURN",

	Variable: "VARIABLE",
	Int:      "VARIABLE",
	Float:    "VARIABLE",
	Float3:   "VARIABLE",
	Float4:   "VARIABLE",
	Bool:     "VARIABLE",

	// Builtins
	Read:      "READ",
	Write:     "WRITE",
	Assert:    "ASSERT",
	Array:     "ARRAY",
	Show:      "SHOW",
	Time:      "TIME",
	Attribute: "ATTRIBUTE",
	To:        "TO",
	Sum:       "SUM",
	Print:     "PRINT",

	// Literals:
	IntLiteral:   "INTVAL",
	FloatLiteral: "FLOATVAL",
	String:       "STRING",
	BoolLiteral:  "VARIABLE",
	Str:          "STR",
}

func (t Token) DumpString() string {
	prefix := typeToDump[t.Type]

	if t.Type == NewLine || t.Type == EOF {
		return prefix
	}

	return fmt.Sprintf("%s '%s'", prefix, t.Val)
}
