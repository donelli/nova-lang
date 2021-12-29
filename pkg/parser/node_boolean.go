package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type BooleanNode struct {
	Value     bool
	nodeRange *shared.Range
}

func NewBooleanNode(value bool, token *lexer.LexerToken) *BooleanNode {
	return &BooleanNode{
		Value:     value,
		nodeRange: token.Range,
	}
}

func (l *BooleanNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *BooleanNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *BooleanNode) Type() ParserNodeType {
	return Node_Bool
}

func (l *BooleanNode) ToHTML() string {
	return BuildNodeBoxHTML("bool", "value-node", fmt.Sprintf("%v", l.Value))
}

func (l *BooleanNode) String() string {
	return fmt.Sprintf("Bool{Value: %v, Range: %v}", l.Value, l.nodeRange)
}

func (l *BooleanNode) Range() *shared.Range {
	return l.nodeRange
}
