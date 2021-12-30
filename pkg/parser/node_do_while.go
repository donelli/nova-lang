package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type DoWhileNode struct {
	Condition Node
	Body      Node
	nodeRange *shared.Range
}

func NewDoWhileNode(condition Node, body Node, startPos shared.Position, endPos shared.Position) *DoWhileNode {
	return &DoWhileNode{
		Condition: condition,
		Body:      body,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *DoWhileNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *DoWhileNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *DoWhileNode) Type() ParserNodeType {
	return Node_DoWhile
}

func (l *DoWhileNode) ToHTML() string {
	return BuildNodeBoxHTML("DO&nbsp;WHILE", "do-while-node", "<div>Condition:<br>"+l.Condition.ToHTML()+"</div><div>Body:<br>"+l.Body.ToHTML()+"</div>")
}

func (l *DoWhileNode) String() string {
	return fmt.Sprintf("DoWhileNode{Cond: %v, Body: %v, Range: %v}", l.Condition, l.Body, l.nodeRange)
}

func (l *DoWhileNode) Range() *shared.Range {
	return l.nodeRange
}
