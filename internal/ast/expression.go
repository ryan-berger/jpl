package ast

type Expression interface {
	Command
	expressionCommand()
}