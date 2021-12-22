package shared

import "fmt"

type Warning struct {
	Message string `json:"message"`
	Range   *Range `json:"range"`
}

func NewWarning(startPos Position, endPos Position, message string) *Warning {
	return &Warning{
		Message: message,
		Range:   NewRange(startPos, endPos),
	}
}

func NewWarningRange(warningRange *Range, message string) *Warning {
	return &Warning{
		Message: message,
		Range:   warningRange,
	}
}

func (w *Warning) StringWithProgram(program string) string {
	return fmt.Sprintf("%s\n  at %s:%v:%v to %s:%v:%v", w.Message, program, w.Range.Start.Row, w.Range.Start.Column, program, w.Range.End.Row, w.Range.End.Column)
}
