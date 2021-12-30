package parser

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/shared"
)

type CommentNode struct {
	Value     string
	nodeRange *shared.Range
}

func NewCommentNode(token *lexer.LexerToken) *CommentNode {
	return &CommentNode{
		Value:     token.Value,
		nodeRange: token.Range,
	}
}

func (l *CommentNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *CommentNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *CommentNode) Type() ParserNodeType {
	return Node_Comment
}

func (l *CommentNode) ToHTML() string {
	return BuildNodeBoxHTML("Com", "comment-node", l.Value)
}

func (l *CommentNode) String() string {
	return fmt.Sprintf("CommentNode{Val: %v, Range: %v}", l.Value, l.nodeRange)
}

func (l *CommentNode) Range() *shared.Range {
	return l.nodeRange
}
