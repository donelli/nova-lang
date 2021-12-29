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
	nodeRange      *shared.Range
}

func NewBinaryOperationNode(leftNode Node, operationToken *lexer.LexerToken, rightNode Node) *BinaryOperationNode {
	return &BinaryOperationNode{
		leftNode:       leftNode,
		operationToken: operationToken,
		rightNode:      rightNode,
		nodeRange:      shared.NewRange(leftNode.StartPos(), rightNode.EndPos()),
	}
}

func (l *BinaryOperationNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *BinaryOperationNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *BinaryOperationNode) Type() ParserNodeType {
	return Node_BinOp
}

func (l *BinaryOperationNode) ToHTML() string {
	return BuildNodeBoxHTML("", "bin-op-node", l.leftNode.ToHTML(), l.operationToken.Value, l.rightNode.ToHTML())
}

func (l *BinaryOperationNode) String() string {
	return fmt.Sprintf("BinOp{Left: %v, Op: %v, Right: %v, Range: %v}", l.leftNode, l.operationToken, l.rightNode, l.nodeRange)
}

func (l *BinaryOperationNode) Range() *shared.Range {
	return l.nodeRange
}
