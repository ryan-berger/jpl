package typed

import (
	"fmt"

	"github.com/ryan-berger/jpl/internal/meta"
)

type Error struct {
	msg string
	meta.Locationer
}

func (e *Error) Error() string {
	line, pos := e.Loc()
	return fmt.Sprintf("%s at position %d:%d", e.msg, line, pos)
}

func NewError(loc meta.Locationer, msg string, args ...interface{}) error {
	return &Error{
		msg:        fmt.Sprintf(msg, args...),
		Locationer: loc,
	}
}
