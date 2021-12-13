package parser

import "recital_lsp/pkg/shared"

type ListNode struct {
	Nodes    []Node
	startPos *shared.Position
	endPos   *shared.Position
}

func NewListNode(Nodes []Node, startPos *shared.Position, endPos *shared.Position) *ListNode {
	return &ListNode{
		Nodes:    Nodes,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *ListNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *ListNode) EndPos() *shared.Position {
	return l.endPos
}
