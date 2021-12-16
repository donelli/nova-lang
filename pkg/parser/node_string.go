package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type StringNode struct {
	Value    string
	startPos *shared.Position
	endPos   *shared.Position
}

func NewStringNode(token *lexer.LexerToken) *StringNode {
	return &StringNode{
		Value:    token.Value,
		startPos: &token.Range.Start,
		endPos:   &token.Range.End,
	}
}

func (l *StringNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *StringNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *StringNode) ToHTML() string {
	return BuildNodeBoxHTML("str", "value-node", fmt.Sprintf("\"%v\"", l.Value))
}

func (l *StringNode) String() string {
	return fmt.Sprintf("String{Value: %v, startPos: %v, endPos: %v}", l.Value, l.startPos, l.endPos)
}
