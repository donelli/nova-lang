package parser

import "recital_lsp/pkg/shared"

type ListNode struct {
	Statements []*StatementNode
	*Node
}

func NewListNode(statements []*StatementNode, startPos *shared.Position, endtPos *shared.Position) *ListNode {
	return &ListNode{
		Statements: statements,
		Node:       NewNode(startPos, endtPos),
	}
}
