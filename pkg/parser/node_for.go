package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type ForNode struct {
	VarName   string
	StartNode Node
	EndNode   Node
	StepNode  Node
	BodyNode  Node
	nodeRange *shared.Range
}

func NewForNode(varName string, startNode Node, endNode Node, stepNode Node, bodyNode Node, startPos shared.Position, endPos shared.Position) *ForNode {
	return &ForNode{
		VarName:   varName,
		StartNode: startNode,
		EndNode:   endNode,
		StepNode:  stepNode,
		BodyNode:  bodyNode,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *ForNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *ForNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *ForNode) Type() ParserNodeType {
	return Node_ForLoop
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

	return fmt.Sprintf("ForNode{Var: %v, Start: %v, End: %v, Step: %v, Range: %v}", l.VarName, l.StartNode, l.EndNode, step, l.nodeRange)
}

func (l *ForNode) Range() *shared.Range {
	return l.nodeRange
}
