package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type NumberNode struct {
	Value     float64
	nodeRange *shared.Range
}

func NewNumberNode(value float64, token *lexer.LexerToken) *NumberNode {
	return &NumberNode{
		Value:     value,
		nodeRange: token.Range,
	}
}

func (l *NumberNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *NumberNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *NumberNode) Type() ParserNodeType {
	return Node_Number
}

func (l *NumberNode) ToHTML() string {
	return BuildNodeBoxHTML("num", "value-node", fmt.Sprintf("%v", l.Value))
}

func (l *NumberNode) String() string {
	return fmt.Sprintf("Number{Value: %v, Range: %v}", l.Value, l.nodeRange)
}

func (l *NumberNode) Range() *shared.Range {
	return l.nodeRange
}
