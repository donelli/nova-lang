package parser

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/shared"
)

type UnaryOperationNode struct {
	OperationToken *lexer.LexerToken
	Node           Node
	nodeRange      *shared.Range
}

func NewUnaryOperationNode(operationToken *lexer.LexerToken, node Node) *UnaryOperationNode {
	return &UnaryOperationNode{
		Node:           node,
		OperationToken: operationToken,
		nodeRange:      shared.NewRange(operationToken.Range.Start, node.EndPos()),
	}
}

func (l *UnaryOperationNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *UnaryOperationNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *UnaryOperationNode) Type() ParserNodeType {
	return Node_UnaryOp
}

func (l *UnaryOperationNode) ToHTML() string {
	return BuildNodeBoxHTML("", "bin-op-node", l.OperationToken.Value, l.Node.ToHTML())
}

func (l *UnaryOperationNode) String() string {
	return fmt.Sprintf("UnaryNode{Oper: %v, Node: %v, Range: %v}", l.OperationToken, l.Node, l.nodeRange)
}

func (l *UnaryOperationNode) Range() *shared.Range {
	return l.nodeRange
}
