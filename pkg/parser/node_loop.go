package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type LoopNode struct {
	startPos *shared.Position
	endPos   *shared.Position
}

func NewLoopNode(tokenRange *shared.Range) *LoopNode {
	return &LoopNode{
		startPos: &tokenRange.Start,
		endPos:   &tokenRange.End,
	}
}

func (l *LoopNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *LoopNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *LoopNode) ToHTML() string {
	return BuildNodeBoxHTML("", "loop-node", "loop")
}

func (l *LoopNode) String() string {
	return fmt.Sprintf("Loop{startPos: %v, endPos: %v}", l.startPos, l.endPos)
}
