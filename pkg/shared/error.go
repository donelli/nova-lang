package shared

import "fmt"

type Error struct {
	Message string
	Range   *Range
	Type    string
}

func NewError(startPos Position, endPos Position, message string) *Error {
	return &Error{
		Message: message,
		Type:    "",
		Range:   NewRange(startPos, endPos),
	}
}

func (e Error) String() string {
	errorType := ""
	if e.Type != "" {
		errorType = e.Type + ": "
	}
	return errorType + e.Message + " at " + e.Range.String()
}

func (e Error) StringWithProgram(program string) string {
	errorType := ""
	if e.Type != "" {
		errorType = e.Type + ": "
	}
	return fmt.Sprintf("%s%s\n  at %s:%v:%v to %s:%v:%v", errorType, e.Message, program, e.Range.Start.Row, e.Range.Start.Column, program, e.Range.End.Row, e.Range.End.Column)
}

func NewInvalidSyntaxErrorRange(errRange *Range, message string) *Error {
	return NewInvalidSyntaxError(errRange.Start, errRange.End, message)
}

func NewInvalidSyntaxError(startPos Position, endPos Position, message string) *Error {
	return &Error{
		Message: message,
		Type:    "Invalid Syntax Error",
		Range:   NewRange(startPos, endPos),
	}
}
