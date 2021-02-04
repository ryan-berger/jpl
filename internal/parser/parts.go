package parser

type Command interface {
	TokenLiteral() string
	String() string
}

type Binding interface {
	String() string
	binding()
}

type Expression interface {
	Command
	expressionCommand()
}

type Statement interface {
	Command
	statementCommand()
}
