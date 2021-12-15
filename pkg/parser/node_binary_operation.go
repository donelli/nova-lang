package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type BinaryOperationNode struct {
	leftNode       Node
	operationToken *lexer.LexerToken
	rightNode      Node
	startPos       *shared.Position
	endPos         *shared.Position
}

func NewBinaryOperationNode(leftNode Node, operationToken *lexer.LexerToken, rightNode Node) *BinaryOperationNode {
	return &BinaryOperationNode{
		leftNode:       leftNode,
		operationToken: operationToken,
		rightNode:      rightNode,
		startPos:       leftNode.StartPos(),
		endPos:         rightNode.EndPos(),
	}
}

func (l *BinaryOperationNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *BinaryOperationNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *BinaryOperationNode) ToHTML() string {
	return BuildNodeBoxHTML("", "bin-op-node", l.leftNode.ToHTML(), l.operationToken.Value, l.rightNode.ToHTML())
}

func (l *BinaryOperationNode) String() string {
	return fmt.Sprintf("BinOp{Left: %v, Op: %v, Right: %v, startPos: %v, endPos: %v}", l.leftNode, l.operationToken, l.rightNode, l.startPos, l.endPos)
}
