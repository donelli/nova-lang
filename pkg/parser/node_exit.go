package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type ExitNode struct {
	startPos *shared.Position
	endPos   *shared.Position
}

func NewExitNode(tokenRange *shared.Range) *ExitNode {
	return &ExitNode{
		startPos: &tokenRange.Start,
		endPos:   &tokenRange.End,
	}
}

func (l *ExitNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *ExitNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *ExitNode) Type() ParserNodeType {
	return Node_Exit
}

func (l *ExitNode) ToHTML() string {
	return BuildNodeBoxHTML("", "exit-node", "exit")
}

func (l *ExitNode) String() string {
	return fmt.Sprintf("Exit{startPos: %v, endPos: %v}", l.startPos, l.endPos)
}
