package parser

import (
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type BinaryOperation struct {
	leftNode       Node
	operationToken *lexer.LexerToken
	rightNode      Node
	startPos       *shared.Position
	endPos         *shared.Position
}

func NewBinaryOperationNode(leftNode Node, operationToken *lexer.LexerToken, rightNode Node) *BinaryOperation {
	return &BinaryOperation{
		leftNode:       leftNode,
		operationToken: operationToken,
		rightNode:      rightNode,
		startPos:       leftNode.StartPos(),
		endPos:         rightNode.EndPos(),
	}
}

func (l *BinaryOperation) StartPos() *shared.Position {
	return l.startPos
}

func (l *BinaryOperation) EndPos() *shared.Position {
	return l.endPos
}
