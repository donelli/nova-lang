package shared

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

func NewInvalidSyntaxError(startPos Position, endPos Position, message string) *Error {
	return &Error{
		Message: message,
		Type:    "Invalid Syntax Error",
		Range:   NewRange(startPos, endPos),
	}
}
