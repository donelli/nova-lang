package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type ClearNode struct {
	Argument string
	Expr     Node
	startPos *shared.Position
	endPos   *shared.Position
}

func NewClearNode(arg string, startPos *shared.Position, endPos *shared.Position) *ClearNode {
	return &ClearNode{
		Argument: arg,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *ClearNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *ClearNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *ClearNode) Type() ParserNodeType {
	return Node_Clear
}

func (l *ClearNode) ToHTML() string {
	return BuildNodeBoxHTML("", "clear-node", "clear", l.Argument)
}

func (l *ClearNode) String() string {
	return fmt.Sprintf("Clear{Arg: %v, startPos: %v, endPos: %v}", l.Argument, l.startPos, l.endPos)
}
