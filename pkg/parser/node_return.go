package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type ReturnNode struct {
	Expr      Node
	ToMaster  bool
	nodeRange *shared.Range
}

func NewReturnNode(expr Node, toMaster bool, startPos shared.Position, endPos shared.Position) *ReturnNode {
	return &ReturnNode{
		Expr:      expr,
		ToMaster:  toMaster,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *ReturnNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *ReturnNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *ReturnNode) Type() ParserNodeType {
	return Node_Return
}

func (l *ReturnNode) ToHTML() string {
	toMaster := "toMaster: " + fmt.Sprintf("%v", l.ToMaster)
	if l.Expr != nil {
		return BuildNodeBoxHTML("", "bin-op-node", "return", toMaster, l.Expr.ToHTML())
	} else {
		return BuildNodeBoxHTML("", "bin-op-node", "return", toMaster)
	}
}

func (l *ReturnNode) String() string {
	return fmt.Sprintf("Return{Expr: %v, ToMaster: %v, Range: %v}", l.Expr, l.ToMaster, l.nodeRange)
}

func (l *ReturnNode) Range() *shared.Range {
	return l.nodeRange
}
