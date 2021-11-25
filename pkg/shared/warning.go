package shared

type Warning struct {
	Message string
	Range   Range
}

func NewWarning(startPos Position, endPos Position, message string) *Warning {
	return &Warning{
		Message: message,
		Range:   *NewRange(startPos, endPos),
	}
}
