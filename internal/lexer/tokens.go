package lexer

type TokenType int

type Token struct {
	Type TokenType
	Val  string
}

const (
	Type TokenType = iota
	Error

	Let
	Function
	Assert

	Identifier

	Boolean
	EOF

	// Organization
	LCurly
	RCurly
	LBrace
	RBrace
	LParen
	RParen
	Comma

	// operators
	Equal
)
