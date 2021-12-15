package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type ReturnNode struct {
	Expr     Node
	startPos *shared.Position
	endPos   *shared.Position
}

func NewReturnNode(expr Node, startPos *shared.Position, endPos *shared.Position) *ReturnNode {
	return &ReturnNode{
		Expr:     expr,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *ReturnNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *ReturnNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *ReturnNode) ToHTML() string {
	if l.Expr != nil {
		return BuildNodeBoxHTML("", "bin-op-node", "return", l.Expr.ToHTML())
	} else {
		return BuildNodeBoxHTML("", "bin-op-node", "return")
	}
}

func (l *ReturnNode) String() string {
	return fmt.Sprintf("Return{Expr: %v, startPos: %v, endPos: %v}", l.Expr, l.startPos, l.endPos)
}
