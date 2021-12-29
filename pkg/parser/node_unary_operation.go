package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type UnaryOperationNode struct {
	operationToken *lexer.LexerToken
	node           Node
	nodeRange      *shared.Range
}

func NewUnaryOperationNode(operationToken *lexer.LexerToken, node Node) *UnaryOperationNode {
	return &UnaryOperationNode{
		node:           node,
		operationToken: operationToken,
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
	return BuildNodeBoxHTML("", "bin-op-node", l.operationToken.Value, l.node.ToHTML())
}

func (l *UnaryOperationNode) String() string {
	return fmt.Sprintf("UnaryNode{Oper: %v, Node: %v, Range: %v}", l.operationToken, l.node, l.nodeRange)
}

func (l *UnaryOperationNode) Range() *shared.Range {
	return l.nodeRange
}
