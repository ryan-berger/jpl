package types

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

type Error struct {
	msg string
	ast.Locationer
}

func (e *Error) Error() string {
	line, pos := e.Loc()
	return fmt.Sprintf("%s at position %d:%d", e.msg, line, pos)
}

func NewError(loc ast.Locationer, msg string, args ...interface{}) error {
	return &Error{
		msg:        fmt.Sprintf(msg, args...),
		Locationer: loc,
	}
}
