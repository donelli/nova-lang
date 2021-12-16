package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type ForNode struct {
	VarName   string
	StartNode Node
	EndNode   Node
	StepNode  Node
	BodyNode  Node
	startPos  *shared.Position
	endPos    *shared.Position
}

func NewForNode(varName string, startNode Node, endNode Node, stepNode Node, bodyNode Node, startPos *shared.Position, endPos *shared.Position) *ForNode {
	return &ForNode{
		VarName:   varName,
		StartNode: startNode,
		EndNode:   endNode,
		StepNode:  stepNode,
		BodyNode:  bodyNode,
		startPos:  startPos,
		endPos:    endPos,
	}
}

func (l *ForNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *ForNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *ForNode) ToHTML() string {

	step := "1"
	if l.StepNode != nil {
		step = l.StepNode.ToHTML()
	}

	return BuildNodeBoxHTML("FOR", "for-node",
		"<div>Variable:<br>"+l.VarName+"</div>"+
			"<div>Start:<br>"+l.StartNode.ToHTML()+"</div>"+
			"<div>End:<br>"+l.EndNode.ToHTML()+"</div>"+
			"<div>Step:<br>"+step+"</div>"+
			"<div>Body:<br>"+l.BodyNode.ToHTML()+"</div>")
}

func (l *ForNode) String() string {

	step := "1"
	if l.StepNode != nil {
		step = fmt.Sprintf("%v", l.StepNode)
	}

	return fmt.Sprintf("ForNode{Var: %v, Start: %v, End: %v, Step: %v, startPos: %v, endPos: %v}", l.VarName, l.StartNode, l.EndNode, step, l.startPos, l.endPos)
}
