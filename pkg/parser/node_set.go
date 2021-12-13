package parser

import "recital_lsp/pkg/shared"

type SetNode struct {
	configName string
	ValueNode  Node
	BoolValue  bool
	startPos   *shared.Position
	endPos     *shared.Position
}

func NewEmptySetNode(configName string, startPos *shared.Position, endPos *shared.Position) *SetNode {
	return &SetNode{
		configName: configName,
		startPos:   startPos,
		endPos:     endPos,
	}
}

func NewSetNode(configName string, valueNode Node, startPos *shared.Position, endPos *shared.Position) *SetNode {
	return &SetNode{
		configName: configName,
		startPos:   startPos,
		endPos:     endPos,
		ValueNode:  valueNode,
	}
}

func NewBoolSetNode(configName string, boolValue bool, startPos *shared.Position, endPos *shared.Position) *SetNode {
	return &SetNode{
		configName: configName,
		startPos:   startPos,
		endPos:     endPos,
		BoolValue:  boolValue,
	}
}

func (l *SetNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *SetNode) EndPos() *shared.Position {
	return l.endPos
}
