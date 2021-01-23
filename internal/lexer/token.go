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

	// Types
	Int
	Float
	Float3
	Float4
	Bool
)

type Token struct {
	Type TokenType
	Val  string
}

var typeToDump = map[TokenType]string{
	ILLEGAL: "ERROR",
	EOF:     "EOF",

	// Delimiters
	LParen:  "LPAREN",
	RParen:  "RPAREN",
	LBrace:  "LBRACE",
	RBrace:  "RBRACE",
	LCurly:  "LCURLY",
	RCurly:  "RCURLY",
	Comma:   "COMMA",
	Colon:   "COLON",
	NewLine: "NEWLINE",

	// Operators
	Assign:             "EQUALS",
	Plus:               "BINOP",
	Minus:              "BINOP",
	Divide:             "BINOP",
	Multiply:           "BIONP",
	Mod:                "BINOP",
	Not:                "BOOLNOT",
	And:                "BINOP",
	Or:                 "BINOP",
	LessThan:           "BINOP",
	LessThanOrEqual:    "BINOP",
	GreaterThan:        "BINOP",
	GreaterThanOrEqual: "BIONP",
	EqualTo:            "BIONP",
	NotEqualTo:         "BIONOP",

	// Keywords
	Function: "FN",
	Let:      "LET",
	If:       "IF",
	Then:     "THEN",
	Else:     "ELSE",
	Return:   "RETURN",

	Variable: "VARIABLE",

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
	IntLiteral: "INTVAL",
	FloatLiteral: "FLOATVAL",
}

func (t *Token) DumpString() string {
	prefix := typeToDump[t.Type]

	if t.Type == NewLine {
		return prefix
	}

	return fmt.Sprintf("%s '%s'", prefix, t.Val)
}
