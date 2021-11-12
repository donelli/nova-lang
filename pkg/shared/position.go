package shared

import (
	"fmt"
)

type Position struct {
	Row    int32
	Column int32
}

func (p Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row, p.Column)
}
