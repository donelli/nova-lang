package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type UnaryOperation struct {
	operationToken *lexer.LexerToken
	node           Node
	startPos       *shared.Position
	endPos         *shared.Position
}

func NewUnaryOperationNode(operationToken *lexer.LexerToken, node Node) *UnaryOperation {
	return &UnaryOperation{
		node:           node,
		operationToken: operationToken,
		startPos:       &operationToken.Range.Start,
		endPos:         node.EndPos(),
	}
}

func (l *UnaryOperation) StartPos() *shared.Position {
	return l.startPos
}

func (l *UnaryOperation) EndPos() *shared.Position {
	return l.endPos
}

func (l *UnaryOperation) ToHTML() string {
	return BuildNodeBoxHTML("", "bin-op-node", l.operationToken.Value, l.node.ToHTML())
}

func (l *UnaryOperation) String() string {
	return fmt.Sprintf("UnaryNode{Oper: %v, Node: %v, startPos: %v, endPos: %v}", l.operationToken, l.node, l.startPos, l.endPos)
}
