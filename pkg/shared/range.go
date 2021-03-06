package shared

import (
	"fmt"
)

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

func (p Range) String() string {
	return fmt.Sprintf("%s-%s", p.Start.String(), p.End.String())
}

func NewRange(startPos Position, endPos Position) *Range {
	return &Range{
		Start: startPos,
		End:   endPos,
	}
}
