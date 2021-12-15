package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type CommentNode struct {
	Value    string
	startPos *shared.Position
	endPos   *shared.Position
}

func NewCommentNode(token *lexer.LexerToken) *CommentNode {
	return &CommentNode{
		Value:    token.Value,
		startPos: &token.Range.Start,
		endPos:   &token.Range.End,
	}
}

func (l *CommentNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *CommentNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *CommentNode) ToHTML() string {
	return BuildNodeBoxHTML("Com", "comment-node", l.Value)
}

func (l *CommentNode) String() string {
	return fmt.Sprintf("CommentNode{Val: %v, startPos: %v, endPos: %v}", l.Value, l.startPos, l.endPos)
}
