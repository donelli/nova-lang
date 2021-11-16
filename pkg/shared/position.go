package shared

import (
	"fmt"
)

type Position struct {
	Row    int32
	Column int32
	Index  int32
}

func (p *Position) String() string {
	return fmt.Sprintf("(%d,%d)", p.Row, p.Column)
}

func (p *Position) Advance(char rune) {

	p.Index++
	p.Column++

	if char == '\n' {
		p.Row++
		p.Column = 0
	}

}

func NewPosition() *Position {
	return &Position{
		Row:    0,
		Column: -1,
		Index:  -1,
	}
}
