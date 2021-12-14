package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type SetNode struct {
	configName string
	ValueNode  Node
	BoolValue  string
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

func NewBoolSetNode(configName string, boolValue string, startPos *shared.Position, endPos *shared.Position) *SetNode {
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

func (l *SetNode) ToHTML() string {
	panic("Not implemented")
}

func (l *SetNode) String() string {

	if l.BoolValue != "" {
		return fmt.Sprintf("SetBoolNode{Config: %v, Bool: %v, startPos: %v, endPos: %v}", l.configName, l.BoolValue, l.startPos, l.endPos)
	}

	if l.ValueNode != nil {
		return fmt.Sprintf("SetNode{Config: %v, Expr: %v, startPos: %v, endPos: %v}", l.configName, l.BoolValue, l.startPos, l.endPos)
	}

	return fmt.Sprintf("SetEmptyNode{Config: %v, startPos: %v, endPos: %v}", l.configName, l.startPos, l.endPos)

}
