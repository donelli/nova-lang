package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

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

func (l *ListNode) ToHTML() string {

	str := "<div class=\"node node-list\">"

	for i := range l.Nodes {
		str += l.Nodes[i].ToHTML() + "<hr>"
	}

	return str + "</div>"
}

func (l *ListNode) String() string {
	return fmt.Sprintf("ListNode{Nodes: %v, startPos: %v, endPos: %v}", l.Nodes, l.startPos, l.endPos)
}
