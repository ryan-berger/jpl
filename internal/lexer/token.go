package lexer

type TokenType int

type Token struct {
	Type TokenType
	Val  string
}

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	NUMBER

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
	Sum
	Print

	// Types
	Int
	Float
	Float3
	Float4
	Bool
)
