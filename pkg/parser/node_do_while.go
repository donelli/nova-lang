package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type DoWhileNode struct {
	Condition Node
	Body      Node
	startPos  *shared.Position
	endPos    *shared.Position
}

func NewDoWhileNode(condition Node, body Node, startPos *shared.Position, endPos *shared.Position) *DoWhileNode {
	return &DoWhileNode{
		Condition: condition,
		Body:      body,
		startPos:  startPos,
		endPos:    endPos,
	}
}

func (l *DoWhileNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *DoWhileNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *DoWhileNode) ToHTML() string {
	return BuildNodeBoxHTML("DO&nbsp;WHILE", "do-while-node", "<div>Condition:<br>"+l.Condition.ToHTML()+"</div><div>Body:<br>"+l.Body.ToHTML()+"</div>")
}

func (l *DoWhileNode) String() string {
	return fmt.Sprintf("DoWhileNode{Cond: %v, Body: %v, startPos: %v, endPos: %v}", l.Condition, l.Body, l.startPos, l.endPos)
}
