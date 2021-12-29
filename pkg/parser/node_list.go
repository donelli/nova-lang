package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type ListNode struct {
	Nodes     []Node
	nodeRange *shared.Range
}

func NewListNode(Nodes []Node, startPos shared.Position, endPos shared.Position) *ListNode {
	return &ListNode{
		Nodes:     Nodes,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *ListNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *ListNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *ListNode) Type() ParserNodeType {
	return Node_List
}

func (l *ListNode) ToHTML() string {

	str := "<div class=\"node node-list\">"

	for i := range l.Nodes {
		str += l.Nodes[i].ToHTML()
		if i < len(l.Nodes)-1 {
			str += "<hr>"
		}
	}

	return str + "</div>"
}

func (l *ListNode) String() string {

	str := ""
	for i := range l.Nodes {
		str += fmt.Sprintf("\n%s", l.Nodes[i])
	}

	return fmt.Sprintf("ListNode{Nodes: %v\n, Range: %v}", str, l.nodeRange)
}

func (l *ListNode) Range() *shared.Range {
	return l.nodeRange
}
