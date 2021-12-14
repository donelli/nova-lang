package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type NumberNode struct {
	Value    float64
	startPos *shared.Position
	endPos   *shared.Position
}

func NewNumberNode(value float64, token *lexer.LexerToken) *NumberNode {
	return &NumberNode{
		Value:    value,
		startPos: &token.Range.Start,
		endPos:   &token.Range.End,
	}
}

func (l *NumberNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *NumberNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *NumberNode) ToHTML() string {
	return fmt.Sprintf("<div class=\"node node-number\">%v</div>", l.Value)
}

func (l *NumberNode) String() string {
	return fmt.Sprintf("Number{Value: %v, startPos: %v, endPos: %v}", l.Value, l.startPos, l.endPos)
}
