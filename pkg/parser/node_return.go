package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type ReturnNode struct {
	Expr     Node
	ToMaster bool
	startPos *shared.Position
	endPos   *shared.Position
}

func NewReturnNode(expr Node, toMaster bool, startPos *shared.Position, endPos *shared.Position) *ReturnNode {
	return &ReturnNode{
		Expr:     expr,
		ToMaster: toMaster,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *ReturnNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *ReturnNode) EndPos() *shared.Position {
	return l.endPos
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
	return fmt.Sprintf("Return{Expr: %v, toMaster: %v, startPos: %v, endPos: %v}", l.Expr, l.ToMaster, l.startPos, l.endPos)
}
