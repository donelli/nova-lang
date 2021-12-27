package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type MacroNode struct {
	Expr     Node
	startPos *shared.Position
	endPos   *shared.Position
}

func NewMacroNode(expr Node, startPos *shared.Position, endPos *shared.Position) *MacroNode {
	return &MacroNode{
		Expr:     expr,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *MacroNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *MacroNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *MacroNode) Type() ParserNodeType {
	return Node_Macro
}

func (l *MacroNode) ToHTML() string {
	return BuildNodeBoxHTML("Macro", "macro-node", l.Expr.ToHTML())
}

func (l *MacroNode) String() string {
	return fmt.Sprintf("MacroNode{Expr: %v, startPos: %v, endPos: %v}", l.Expr, l.startPos, l.endPos)
}
