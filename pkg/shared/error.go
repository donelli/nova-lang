package shared

type Error struct {
	Message string
	Range   Range
}

func NewError(startPos Position, endPos Position, message string) Error {
	return Error{
		Message: message,
		Range:   *NewRange(startPos, endPos),
	}
}

func (e Error) String() string {
	return e.Message + " at " + e.Range.String()
}
