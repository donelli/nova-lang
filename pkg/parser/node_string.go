package parser

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/shared"
)

type StringNode struct {
	Value     string
	nodeRange *shared.Range
}

func NewStringNode(token *lexer.LexerToken) *StringNode {
	return &StringNode{
		Value:     token.Value,
		nodeRange: token.Range,
	}
}

func (l *StringNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *StringNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *StringNode) Type() ParserNodeType {
	return Node_String
}

func (l *StringNode) ToHTML() string {
	return BuildNodeBoxHTML("str", "value-node", fmt.Sprintf("\"%v\"", l.Value))
}

func (l *StringNode) String() string {
	return fmt.Sprintf("String{Value: %v, Range: %v}", l.Value, l.nodeRange)
}

func (l *StringNode) Range() *shared.Range {
	return l.nodeRange
}
