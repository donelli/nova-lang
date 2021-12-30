package parser

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/shared"
)

type BinaryOperationNode struct {
	LeftNode       Node
	OperationToken *lexer.LexerToken
	RightNode      Node
	nodeRange      *shared.Range
}

func NewBinaryOperationNode(leftNode Node, operationToken *lexer.LexerToken, rightNode Node) *BinaryOperationNode {
	return &BinaryOperationNode{
		LeftNode:       leftNode,
		OperationToken: operationToken,
		RightNode:      rightNode,
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
	return BuildNodeBoxHTML("", "bin-op-node", l.LeftNode.ToHTML(), l.OperationToken.Value, l.RightNode.ToHTML())
}

func (l *BinaryOperationNode) String() string {
	return fmt.Sprintf("BinOp{Left: %v, Op: %v, Right: %v, Range: %v}", l.LeftNode, l.OperationToken, l.RightNode, l.nodeRange)
}

func (l *BinaryOperationNode) Range() *shared.Range {
	return l.nodeRange
}
