package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type SetNode struct {
	configName string
	ValueNode  Node
	BoolValue  string
	nodeRange  *shared.Range
}

func NewEmptySetNode(configName string, startPos shared.Position, endPos shared.Position) *SetNode {
	return &SetNode{
		configName: configName,
		nodeRange:  shared.NewRange(startPos, endPos),
	}
}

func NewSetNode(configName string, valueNode Node, startPos shared.Position, endPos shared.Position) *SetNode {
	return &SetNode{
		configName: configName,
		nodeRange:  shared.NewRange(startPos, endPos),
		ValueNode:  valueNode,
	}
}

func NewBoolSetNode(configName string, boolValue string, startPos shared.Position, endPos shared.Position) *SetNode {
	return &SetNode{
		configName: configName,
		nodeRange:  shared.NewRange(startPos, endPos),
		BoolValue:  boolValue,
	}
}

func (l *SetNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *SetNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *SetNode) Type() ParserNodeType {
	return Node_Set
}

func (l *SetNode) ToHTML() string {
	panic("Not implemented")
}

func (l *SetNode) String() string {

	if l.BoolValue != "" {
		return fmt.Sprintf("SetBoolNode{Config: %v, Bool: %v, Range: %v}", l.configName, l.BoolValue, l.nodeRange)
	}

	if l.ValueNode != nil {
		return fmt.Sprintf("SetNode{Config: %v, Expr: %v, Range: %v}", l.configName, l.BoolValue, l.nodeRange)
	}

	return fmt.Sprintf("SetEmptyNode{Config: %v, Range: %v}", l.configName, l.nodeRange)

}

func (l *SetNode) Range() *shared.Range {
	return l.nodeRange
}
