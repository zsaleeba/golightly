package golightly

import "fmt"

type Error struct {
	filename string
	pos      SrcSpan
	message  string
}

func NewError(filename string, pos SrcSpan, message string) *Error {
	e := new(Error)
	e.filename = filename
	e.pos = pos
	e.message = message

	return e
}

func (e *Error) Error() string {
	return fmt.Sprint(e.filename, ":", e.pos.start.Line, ": ", e.message)
}
