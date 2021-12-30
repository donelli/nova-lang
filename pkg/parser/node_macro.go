package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type MacroNode struct {
	Expr      Node
	nodeRange *shared.Range
}

func NewMacroNode(expr Node, startPos shared.Position, endPos shared.Position) *MacroNode {
	return &MacroNode{
		Expr:      expr,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *MacroNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *MacroNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *MacroNode) Type() ParserNodeType {
	return Node_Macro
}

func (l *MacroNode) ToHTML() string {
	return BuildNodeBoxHTML("Macro", "macro-node", l.Expr.ToHTML())
}

func (l *MacroNode) String() string {
	return fmt.Sprintf("MacroNode{Expr: %v, Range: %v}", l.Expr, l.nodeRange)
}

func (l *MacroNode) Range() *shared.Range {
	return l.nodeRange
}
