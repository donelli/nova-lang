package shared

import (
	"fmt"
)

type Position struct {
	Row    int32 `json:"row"`
	Column int32 `json:"column"`
	Index  int32 `json:"-"`
}

func (p *Position) String() string {
	return fmt.Sprintf("{%d:%d}", p.Row, p.Column)
}

func (p *Position) Advance(char rune) {

	p.Index++
	p.Column++

	if char == '\n' {
		p.Row++
		p.Column = 0
	}

}

func (p *Position) Copy() *Position {
	return &Position{
		Row:    p.Row,
		Column: p.Column,
		Index:  p.Index,
	}
}

func NewPosition() *Position {
	return &Position{
		Row:    0,
		Column: 0,
		Index:  -1,
	}
}
