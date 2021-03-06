package shared

import "fmt"

type Error struct {
	Message string `json:"message"`
	Range   *Range `json:"range"`
	Type    string `json:"type"`
}

func NewError(startPos Position, endPos Position, message string) *Error {
	return &Error{
		Message: message,
		Type:    "",
		Range:   NewRange(startPos, endPos),
	}
}

func (e *Error) UpdateRange(errorRange *Range) *Error {
	e.Range = errorRange
	return e
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
	return fmt.Sprintf("%s%s\n  at %s:%v:%v to %s:%v:%v", errorType, e.Message, program, e.Range.Start.Row+1, e.Range.Start.Column+1, program, e.Range.End.Row+1, e.Range.End.Column+1)
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

func NewRuntimeError(startPos Position, endPos Position, message string) *Error {
	return &Error{
		Message: message,
		Type:    "Runtime Error",
		Range:   NewRange(startPos, endPos),
	}
}

func NewRuntimeErrorRange(errRange *Range, message string) *Error {
	return NewRuntimeError(errRange.Start, errRange.End, message)
}

func NewAssertError(errRange *Range, message string) *Error {
	return &Error{
		Message: message,
		Type:    "Assert error",
		Range:   errRange,
	}
}
