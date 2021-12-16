package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type BooleanNode struct {
	Value    bool
	startPos *shared.Position
	endPos   *shared.Position
}

func NewBooleanNode(value bool, token *lexer.LexerToken) *BooleanNode {
	return &BooleanNode{
		Value:    value,
		startPos: &token.Range.Start,
		endPos:   &token.Range.End,
	}
}

func (l *BooleanNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *BooleanNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *BooleanNode) Type() ParserNodeType {
	return Node_Bool
}

func (l *BooleanNode) ToHTML() string {
	return BuildNodeBoxHTML("bool", "value-node", fmt.Sprintf("%v", l.Value))
}

func (l *BooleanNode) String() string {
	return fmt.Sprintf("Bool{Value: %v, startPos: %v, endPos: %v}", l.Value, l.startPos, l.endPos)
}
