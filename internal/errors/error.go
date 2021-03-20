package errors

import (
	"fmt"


	"github.com/ryan-berger/jpl/internal/meta"
)

type Type struct {
	Message string
	meta.Locationer
}

func (t *Type) Error() string {
	line, pos := t.Loc()
	return fmt.Sprintf("type error: %s at %d:%d", t.Message, line, pos)
}

func TypeError(message string, loc meta.Locationer) *Type {
	return &Type{
		Message:    message,
		Locationer: loc,
	}
}

type Parse struct {
	Message string
	meta.Locationer
}

func (t *Parse) Error() string {
	line, pos := t.Loc()
	return fmt.Sprintf("parse error: %s at %d:%d", t.Message, line, pos)
}

func ParseError(message string, loc meta.Locationer) *Parse {
	return &Parse{
		Message:    message,
		Locationer: loc,
	}
}
